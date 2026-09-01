package swt3

import (
	"strings"
	"strconv"
)

// Witness provides named witness methods for minting SWT3 anchors.
//
// Each method returns a WitnessPayload with computed factors, fingerprint,
// and optional context (stripped at clearing level >= 2). Payloads are
// accumulated in an internal buffer for batch ingestion.
type Witness struct {
	Config  WitnessConfig
	Buffer  []WitnessPayload
}

// NewWitness creates a Witness with the given configuration.
func NewWitness(config WitnessConfig) *Witness {
	if config.BufferSize == 0 {
		config.BufferSize = 100
	}
	return &Witness{
		Config: config,
		Buffer: make([]WitnessPayload, 0, config.BufferSize),
	}
}

// mint creates a WitnessPayload and enqueues it.
func (w *Witness) mint(procedureID string, fa, fb, fc float64) WitnessPayload {
	ts, epoch := TimestampMs()
	fp := MintFingerprint(w.Config.TenantID, procedureID, fa, fb, fc, ts)
	p := WitnessPayload{
		ProcedureID:          procedureID,
		FactorA:              fa,
		FactorB:              fb,
		FactorC:              fc,
		ClearingLevel:        w.Config.ClearingLevel,
		AnchorFingerprint:    fp,
		AnchorEpoch:          epoch,
		FingerprintTimestamp: ts,
	}
	if w.Config.AgentID != "" {
		p.AgentID = w.Config.AgentID
	}
	if w.Config.Jurisdiction != "" {
		p.Jurisdiction = w.Config.Jurisdiction
	}
	if w.Config.LegalBasis != "" {
		p.LegalBasis = w.Config.LegalBasis
	}
	if w.Config.PurposeClass != "" {
		p.PurposeClass = w.Config.PurposeClass
	}
	if w.Config.SigningKey != "" {
		p.PayloadSignature = SignPayload(w.Config.SigningKey, fp, w.Config.AgentID)
		p.SigningAlgorithm = "hmac-sha256"
	}
	return p
}

// enqueue appends a payload to the buffer.
func (w *Witness) enqueue(p WitnessPayload) {
	w.Buffer = append(w.Buffer, p)
}

// Flush returns all buffered payloads and resets the buffer.
func (w *Witness) Flush() []WitnessPayload {
	out := w.Buffer
	w.Buffer = make([]WitnessPayload, 0, w.Config.BufferSize)
	return out
}

// hashToFactor converts a hex hash string to a float64 factor value.
func hashToFactor(hexStr string) float64 {
	if len(hexStr) > 8 {
		hexStr = hexStr[:8]
	}
	v, _ := strconv.ParseUint(hexStr, 16, 64)
	return float64(v)
}

// ── NHI: Non-Human Identity Governance ──────────────────────────────

// WitnessNhiScope witnesses credential scope attestation (NHI-SCOPE.1).
func (w *Witness) WitnessNhiScope(credentialID, scope string, ttlSeconds int) WitnessPayload {
	credHash := Sha256Truncated(credentialID, 16)
	scopeHash := Sha256Truncated(strings.ToLower(scope), 16)
	fa := hashToFactor(credHash)
	fb := hashToFactor(scopeHash)
	fc := float64(ttlSeconds)
	p := w.mint("NHI-SCOPE.1", fa, fb, fc)
	if w.Config.ClearingLevel <= ClearingStandard {
		p.AIModelID = "nhi-scope-" + credHash[:8]
		p.AIContext = map[string]interface{}{
			"provider":          "nhi-governance",
			"credential_id_hash": credHash,
			"scope_hash":         scopeHash,
			"ttl_seconds":        ttlSeconds,
		}
	}
	w.enqueue(p)
	return p
}

// WitnessNhiLifecycle witnesses a credential lifecycle event (NHI-CYCLE.1).
func (w *Witness) WitnessNhiLifecycle(eventType int, credentialID, issuer string) WitnessPayload {
	credHash := Sha256Truncated(credentialID, 16)
	issuerHash := Sha256Truncated(issuer, 16)
	fa := float64(eventType)
	fb := hashToFactor(credHash)
	fc := hashToFactor(issuerHash)
	p := w.mint("NHI-CYCLE.1", fa, fb, fc)
	if w.Config.ClearingLevel <= ClearingStandard {
		p.AIModelID = "nhi-lifecycle-" + strconv.Itoa(eventType)
		p.AIContext = map[string]interface{}{
			"provider":          "nhi-governance",
			"event_type":         eventType,
			"credential_id_hash": credHash,
			"issuer_hash":        issuerHash,
		}
	}
	w.enqueue(p)
	return p
}

// WitnessNhiPrivilegeChange witnesses a privilege change (NHI-PRIV.1).
func (w *Witness) WitnessNhiPrivilegeChange(credentialID, previousScope, newScope string) WitnessPayload {
	credHash := Sha256Truncated(credentialID, 16)
	fa := hashToFactor(credHash)
	var fb float64
	if previousScope != "" {
		prevHash := Sha256Truncated(strings.ToLower(previousScope), 16)
		fb = hashToFactor(prevHash)
	}
	newHash := Sha256Truncated(strings.ToLower(newScope), 16)
	fc := hashToFactor(newHash)
	p := w.mint("NHI-PRIV.1", fa, fb, fc)
	if w.Config.ClearingLevel <= ClearingStandard {
		prevH := "0"
		if previousScope != "" {
			prevH = Sha256Truncated(strings.ToLower(previousScope), 16)
		}
		p.AIModelID = "nhi-priv-" + credHash[:8]
		p.AIContext = map[string]interface{}{
			"provider":            "nhi-governance",
			"credential_id_hash":  credHash,
			"previous_scope_hash": prevH,
			"new_scope_hash":      Sha256Truncated(strings.ToLower(newScope), 16),
		}
	}
	w.enqueue(p)
	return p
}

// WitnessNhiRotation witnesses credential rotation (NHI-ROTATE.1).
func (w *Witness) WitnessNhiRotation(oldCredentialID, newCredentialID string, reason int) WitnessPayload {
	oldHash := Sha256Truncated(oldCredentialID, 16)
	newHash := Sha256Truncated(newCredentialID, 16)
	fa := hashToFactor(oldHash)
	fb := hashToFactor(newHash)
	fc := float64(reason)
	p := w.mint("NHI-ROTATE.1", fa, fb, fc)
	if w.Config.ClearingLevel <= ClearingStandard {
		p.AIModelID = "nhi-rotate-" + strconv.Itoa(reason)
		p.AIContext = map[string]interface{}{
			"provider":            "nhi-governance",
			"old_credential_hash": oldHash,
			"new_credential_hash": newHash,
			"rotation_reason":     reason,
		}
	}
	w.enqueue(p)
	return p
}

// WitnessNhiDelegation witnesses agent-to-agent credential delegation (NHI-AGENT.1).
func (w *Witness) WitnessNhiDelegation(delegatorCredentialID, delegateeCredentialID string, delegationDepth int) WitnessPayload {
	delegatorHash := Sha256Truncated(delegatorCredentialID, 16)
	delegateeHash := Sha256Truncated(delegateeCredentialID, 16)
	fa := hashToFactor(delegatorHash)
	fb := hashToFactor(delegateeHash)
	depth := delegationDepth
	if depth < 1 {
		depth = 1
	}
	fc := float64(depth)
	p := w.mint("NHI-AGENT.1", fa, fb, fc)
	if w.Config.ClearingLevel <= ClearingStandard {
		p.AIModelID = "nhi-delegation-depth" + strconv.Itoa(depth)
		p.AIContext = map[string]interface{}{
			"provider":          "nhi-governance",
			"delegator_hash":    delegatorHash,
			"delegatee_hash":    delegateeHash,
			"delegation_depth":  depth,
		}
	}
	w.enqueue(p)
	return p
}

// WitnessNhiRevocation witnesses credential revocation (NHI-REVOKE.1).
func (w *Witness) WitnessNhiRevocation(credentialID string, reason int, cascade bool) WitnessPayload {
	credHash := Sha256Truncated(credentialID, 16)
	fa := hashToFactor(credHash)
	fb := float64(reason)
	var fc float64
	if cascade {
		fc = 1
	}
	p := w.mint("NHI-REVOKE.1", fa, fb, fc)
	if w.Config.ClearingLevel <= ClearingStandard {
		p.AIModelID = "nhi-revoke-" + strconv.Itoa(reason)
		p.AIContext = map[string]interface{}{
			"provider":               "nhi-governance",
			"revoked_credential_hash": credHash,
			"reason_code":             reason,
			"cascade":                 cascade,
		}
	}
	w.enqueue(p)
	return p
}

// ── HBOM: Hardware Bill of Materials ────────────────────────────────

// WitnessHardwareInventory witnesses hardware inventory attestation (HBOM-INV.1).
func (w *Witness) WitnessHardwareInventory(componentCount int, manifestHash string, deltaFromBaseline int) WitnessPayload {
	fa := float64(componentCount)
	fb := hashToFactor(manifestHash)
	fc := float64(deltaFromBaseline)
	p := w.mint("HBOM-INV.1", fa, fb, fc)
	if w.Config.ClearingLevel <= ClearingStandard {
		mh := manifestHash
		if len(mh) > 16 {
			mh = mh[:16]
		}
		p.AIModelID = "hbom-inv-" + strconv.Itoa(componentCount)
		p.AIContext = map[string]interface{}{
			"provider":            "hbom-governance",
			"component_count":     componentCount,
			"manifest_hash":       mh,
			"delta_from_baseline": deltaFromBaseline,
		}
	}
	w.enqueue(p)
	return p
}

// WitnessComponentLifecycle witnesses a component lifecycle event (HBOM-LIFE.1).
func (w *Witness) WitnessComponentLifecycle(eventType int, componentID string, ageDays int) WitnessPayload {
	compHash := Sha256Truncated(componentID, 16)
	fa := float64(eventType)
	fb := hashToFactor(compHash)
	fc := float64(ageDays)
	p := w.mint("HBOM-LIFE.1", fa, fb, fc)
	if w.Config.ClearingLevel <= ClearingStandard {
		p.AIModelID = "hbom-lifecycle-" + strconv.Itoa(eventType)
		p.AIContext = map[string]interface{}{
			"provider":       "hbom-governance",
			"event_type":     eventType,
			"component_hash": compHash,
			"age_days":       ageDays,
		}
	}
	w.enqueue(p)
	return p
}

// WitnessThermalProfile witnesses thermal profile attestation (HBOM-THERM.1).
func (w *Witness) WitnessThermalProfile(ambientTempC, componentTempC float64, thresholdExceeded bool) WitnessPayload {
	fa := ambientTempC
	fb := componentTempC
	var fc float64
	if thresholdExceeded {
		fc = 1
	}
	p := w.mint("HBOM-THERM.1", fa, fb, fc)
	if w.Config.ClearingLevel <= ClearingStandard {
		label := "normal"
		if thresholdExceeded {
			label = "alarm"
		}
		p.AIModelID = "hbom-thermal-" + label
		p.AIContext = map[string]interface{}{
			"provider":           "hbom-governance",
			"ambient_temp_c":     ambientTempC,
			"component_temp_c":   componentTempC,
			"threshold_exceeded": thresholdExceeded,
		}
	}
	w.enqueue(p)
	return p
}

// WitnessWaterConsumption witnesses water consumption (HBOM-WATER.1).
func (w *Witness) WitnessWaterConsumption(litersConsumed float64, wueRatioX1000 int, sourceType int) WitnessPayload {
	fa := litersConsumed
	fb := float64(wueRatioX1000)
	fc := float64(sourceType)
	p := w.mint("HBOM-WATER.1", fa, fb, fc)
	if w.Config.ClearingLevel <= ClearingStandard {
		p.AIModelID = "hbom-water-" + strconv.Itoa(sourceType)
		p.AIContext = map[string]interface{}{
			"provider":        "hbom-governance",
			"liters_consumed": litersConsumed,
			"wue_ratio_x1000": wueRatioX1000,
			"source_type":     sourceType,
		}
	}
	w.enqueue(p)
	return p
}

// WitnessPowerUsage witnesses PUE attestation (HBOM-PUE.1).
func (w *Witness) WitnessPowerUsage(totalFacilityKw, itLoadKw float64, pueX1000 int) WitnessPayload {
	fa := totalFacilityKw
	fb := itLoadKw
	fc := float64(pueX1000)
	p := w.mint("HBOM-PUE.1", fa, fb, fc)
	if w.Config.ClearingLevel <= ClearingStandard {
		p.AIModelID = "hbom-pue-" + strconv.Itoa(pueX1000)
		p.AIContext = map[string]interface{}{
			"provider":         "hbom-governance",
			"total_facility_kw": totalFacilityKw,
			"it_load_kw":        itLoadKw,
			"pue_x1000":         pueX1000,
		}
	}
	w.enqueue(p)
	return p
}

// WitnessSupplyChainProvenance witnesses hardware supply chain provenance (HBOM-SUPPLY.1).
func (w *Witness) WitnessSupplyChainProvenance(supplierID string, provenanceVerified bool, countryOfOrigin string) WitnessPayload {
	supplierHash := Sha256Truncated(supplierID, 16)
	countryHash := Sha256Truncated(strings.ToUpper(countryOfOrigin), 16)
	fa := hashToFactor(supplierHash)
	var fb float64
	if provenanceVerified {
		fb = 1
	}
	fc := hashToFactor(countryHash)
	p := w.mint("HBOM-SUPPLY.1", fa, fb, fc)
	if w.Config.ClearingLevel <= ClearingStandard {
		label := "unverified"
		if provenanceVerified {
			label = "verified"
		}
		p.AIModelID = "hbom-supply-" + label
		p.AIContext = map[string]interface{}{
			"provider":             "hbom-governance",
			"supplier_hash":        supplierHash,
			"provenance_verified":  provenanceVerified,
			"country_hash":         countryHash,
		}
	}
	w.enqueue(p)
	return p
}

// ── DPP: Digital Product Passport ───────────────────────────────────

// WitnessBatterySoh witnesses battery state of health (DPP-SOH.1).
func (w *Witness) WitnessBatterySoh(sohPercent float64, cycleCount int, capacityKwh float64) WitnessPayload {
	fa := float64(int(sohPercent * 100))
	fb := float64(cycleCount)
	fc := float64(int(capacityKwh * 100))
	p := w.mint("DPP-SOH.1", fa, fb, fc)
	if w.Config.ClearingLevel <= ClearingStandard {
		p.AIModelID = "dpp-soh-" + strconv.Itoa(int(sohPercent)) + "pct"
		p.AIContext = map[string]interface{}{
			"provider":     "dpp-governance",
			"soh_percent":  sohPercent,
			"cycle_count":  cycleCount,
			"capacity_kwh": capacityKwh,
		}
	}
	w.enqueue(p)
	return p
}

// WitnessChargeCycle witnesses a charge/discharge cycle event (DPP-CHRG.1).
func (w *Witness) WitnessChargeCycle(eventType int, energyKwh, peakTempC float64) WitnessPayload {
	fa := float64(eventType)
	fb := float64(int(energyKwh * 100))
	fc := peakTempC
	p := w.mint("DPP-CHRG.1", fa, fb, fc)
	if w.Config.ClearingLevel <= ClearingStandard {
		p.AIModelID = "dpp-chrg-" + strconv.Itoa(eventType)
		p.AIContext = map[string]interface{}{
			"provider":    "dpp-governance",
			"event_type":  eventType,
			"energy_kwh":  energyKwh,
			"peak_temp_c": peakTempC,
		}
	}
	w.enqueue(p)
	return p
}

// WitnessDegradationEvent witnesses a battery degradation event (DPP-DEGRAD.1).
func (w *Witness) WitnessDegradationEvent(degradationType int, sohDeltaPercent, ambientTempC float64) WitnessPayload {
	fa := float64(degradationType)
	fb := float64(int(sohDeltaPercent * 100))
	fc := ambientTempC
	p := w.mint("DPP-DEGRAD.1", fa, fb, fc)
	if w.Config.ClearingLevel <= ClearingStandard {
		p.AIModelID = "dpp-degrad-" + strconv.Itoa(degradationType)
		p.AIContext = map[string]interface{}{
			"provider":          "dpp-governance",
			"degradation_type":  degradationType,
			"soh_delta_percent": sohDeltaPercent,
			"ambient_temp_c":    ambientTempC,
		}
	}
	w.enqueue(p)
	return p
}

// WitnessEndOfLife witnesses end-of-life / recycling handoff (DPP-EOL.1).
func (w *Witness) WitnessEndOfLife(dispositionType int, handlerID string, finalSohPercent float64) WitnessPayload {
	handlerHash := Sha256Truncated(handlerID, 16)
	fa := float64(dispositionType)
	fb := hashToFactor(handlerHash)
	fc := float64(int(finalSohPercent * 100))
	p := w.mint("DPP-EOL.1", fa, fb, fc)
	if w.Config.ClearingLevel <= ClearingStandard {
		p.AIModelID = "dpp-eol-" + strconv.Itoa(dispositionType)
		p.AIContext = map[string]interface{}{
			"provider":          "dpp-governance",
			"disposition_type":  dispositionType,
			"handler_hash":      handlerHash,
			"final_soh_percent": finalSohPercent,
		}
	}
	w.enqueue(p)
	return p
}

// ── ADR: Automated Demand Response ──────────────────────────────────

// WitnessDemandResponse witnesses a demand response event (ADR-EVENT.1).
func (w *Witness) WitnessDemandResponse(eventPhase int, committedKw float64, signalSource string) WitnessPayload {
	signalHash := Sha256Truncated(signalSource, 16)
	fa := float64(eventPhase)
	fb := committedKw
	fc := hashToFactor(signalHash)
	p := w.mint("ADR-EVENT.1", fa, fb, fc)
	if w.Config.ClearingLevel <= ClearingStandard {
		p.AIModelID = "adr-event-" + strconv.Itoa(eventPhase)
		p.AIContext = map[string]interface{}{
			"provider":          "adr-governance",
			"event_phase":       eventPhase,
			"committed_kw":      committedKw,
			"signal_source_hash": signalHash,
		}
	}
	w.enqueue(p)
	return p
}

// WitnessBaselineConsumption witnesses baseline consumption attestation (ADR-BASE.1).
func (w *Witness) WitnessBaselineConsumption(baselineKw float64, measurementMethod, confidenceX1000 int) WitnessPayload {
	fa := baselineKw
	fb := float64(measurementMethod)
	fc := float64(confidenceX1000)
	p := w.mint("ADR-BASE.1", fa, fb, fc)
	if w.Config.ClearingLevel <= ClearingStandard {
		p.AIModelID = "adr-baseline-" + strconv.Itoa(measurementMethod)
		p.AIContext = map[string]interface{}{
			"provider":           "adr-governance",
			"baseline_kw":        baselineKw,
			"measurement_method": measurementMethod,
			"confidence_x1000":   confidenceX1000,
		}
	}
	w.enqueue(p)
	return p
}

// WitnessCurtailment witnesses curtailment verification (ADR-CURT.1).
func (w *Witness) WitnessCurtailment(actualReductionKw, committedKw float64, complianceRatioX1000 int) WitnessPayload {
	fa := actualReductionKw
	fb := committedKw
	fc := float64(complianceRatioX1000)
	p := w.mint("ADR-CURT.1", fa, fb, fc)
	if w.Config.ClearingLevel <= ClearingStandard {
		p.AIModelID = "adr-curtailment-" + strconv.Itoa(complianceRatioX1000)
		p.AIContext = map[string]interface{}{
			"provider":               "adr-governance",
			"actual_reduction_kw":    actualReductionKw,
			"committed_kw":           committedKw,
			"compliance_ratio_x1000": complianceRatioX1000,
		}
	}
	w.enqueue(p)
	return p
}

// WitnessSettlement witnesses settlement data attestation (ADR-SETTLE.1).
func (w *Witness) WitnessSettlement(settlementKwh, priceUsdPerMwh float64, eventCount int) WitnessPayload {
	fa := float64(int(settlementKwh * 100))
	fb := float64(int(priceUsdPerMwh * 100))
	fc := float64(eventCount)
	p := w.mint("ADR-SETTLE.1", fa, fb, fc)
	if w.Config.ClearingLevel <= ClearingStandard {
		p.AIModelID = "adr-settlement-" + strconv.Itoa(eventCount) + "events"
		p.AIContext = map[string]interface{}{
			"provider":          "adr-governance",
			"settlement_kwh":    settlementKwh,
			"price_usd_per_mwh": priceUsdPerMwh,
			"event_count":       eventCount,
		}
	}
	w.enqueue(p)
	return p
}

// WitnessCarbonCredit witnesses carbon credit / REC attestation (ADR-CARBON.1).
func (w *Witness) WitnessCarbonCredit(creditType int, quantityMwh float64, registryID string) WitnessPayload {
	registryHash := Sha256Truncated(registryID, 16)
	fa := float64(creditType)
	fb := float64(int(quantityMwh * 100))
	fc := hashToFactor(registryHash)
	p := w.mint("ADR-CARBON.1", fa, fb, fc)
	if w.Config.ClearingLevel <= ClearingStandard {
		p.AIModelID = "adr-carbon-" + strconv.Itoa(creditType)
		p.AIContext = map[string]interface{}{
			"provider":      "adr-governance",
			"credit_type":   creditType,
			"quantity_mwh":  quantityMwh,
			"registry_hash": registryHash,
		}
	}
	w.enqueue(p)
	return p
}

// WitnessGridSignal witnesses grid signal correlation (ADR-GRID.1).
func (w *Witness) WitnessGridSignal(signalType int, responseLatencyMs int, gridOperator string) WitnessPayload {
	operatorHash := Sha256Truncated(gridOperator, 16)
	fa := float64(signalType)
	fb := float64(responseLatencyMs)
	fc := hashToFactor(operatorHash)
	p := w.mint("ADR-GRID.1", fa, fb, fc)
	if w.Config.ClearingLevel <= ClearingStandard {
		p.AIModelID = "adr-grid-" + strconv.Itoa(signalType)
		p.AIContext = map[string]interface{}{
			"provider":           "adr-governance",
			"signal_type":        signalType,
			"response_latency_ms": responseLatencyMs,
			"grid_operator_hash":  operatorHash,
		}
	}
	w.enqueue(p)
	return p
}
