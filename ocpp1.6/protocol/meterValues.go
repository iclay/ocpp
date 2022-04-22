package protocol

import (
	validator "github.com/go-playground/validator/v10"
)

func init() {
	Validate.RegisterValidation("readingContext", func(f validator.FieldLevel) bool {
		readingContext := ReadingContext(f.Field().String())
		switch readingContext {
		case interruptionBegin, interruptionEnd, sampleClock, samplePeriodic,
			transcationBegin, transcationEnd, trigger, other:
			return true
		default:
			return false
		}
	})
	Validate.RegisterValidation("valueFormat", func(fl validator.FieldLevel) bool {
		valueFormat := ValueFormat(fl.Field().String())
		switch valueFormat {
		case rawFormat, signFormat:
			return true
		default:
			return false
		}
	})

	Validate.RegisterValidation("measurand", func(f validator.FieldLevel) bool {
		measurand := Measurand(f.Field().String())
		switch measurand {
		case energyActiveExportRegister, energyActiveImportRegister, energyReactiveExportRegister, energyReactiveImportRegister,
			energyActiveExportInterval, energyActiveImportInterval, energyReactiveExportInterval, energyReactiveImportInterval,
			powerActiveExport, powerActiveImport, powerOffered, powerReactiveExport,
			powerReactiveImport, powerFactor, currentImport, currentExport,
			currentOffered, voltage, frequency, temperature, soc, rpm:
			return true
		default:
			return false
		}
	})

	Validate.RegisterValidation("phase", func(f validator.FieldLevel) bool {
		phase := Phase(f.Field().String())
		switch phase {
		case L1, L2, L3, N, L1N, L2N, L3N, L1L2, L2L3, L3L1:
			return true
		default:
			return false
		}
	})

	Validate.RegisterValidation("location", func(f validator.FieldLevel) bool {
		location := Location(f.Field().String())
		switch location {
		case cable, ev, inlet, outlet, body:
			return true
		default:
			return false
		}
	})

	Validate.RegisterValidation("unitOfMeasure", func(f validator.FieldLevel) bool {
		unitOfMeasure := UnitOfMeasure(f.Field().String())
		switch unitOfMeasure {
		case Wh, kWh, varh, kvarh,
			W, kW, VA, kVA,
			Var, kvar, A, V, K,
			Celcius, Celsius, Fahrenheit, Percent:
			return true
		default:
			return false
		}
	})
}

type UnitOfMeasure string

const (
	Wh         UnitOfMeasure = "Wh"
	kWh        UnitOfMeasure = "kWh"
	varh       UnitOfMeasure = "varh"
	kvarh      UnitOfMeasure = "kvarh"
	W          UnitOfMeasure = "W"
	kW         UnitOfMeasure = "kW"
	VA         UnitOfMeasure = "VA"
	kVA        UnitOfMeasure = "kVA"
	Var        UnitOfMeasure = "var"
	kvar       UnitOfMeasure = "kvar"
	A          UnitOfMeasure = "A"
	V          UnitOfMeasure = "V"
	K          UnitOfMeasure = "K"
	Celcius    UnitOfMeasure = "Celcius"
	Celsius    UnitOfMeasure = "Celsius"
	Fahrenheit UnitOfMeasure = "Fahrenheit"
	Percent    UnitOfMeasure = "Percent"
)

type Location string

const (
	cable  Location = "Cable"
	ev     Location = "EV"
	inlet  Location = "Inlet"
	outlet Location = "Outlet"
	body   Location = "Body"
)

type Phase string

const (
	L1   Phase = "L1"
	L2   Phase = "L2"
	L3   Phase = "L3"
	N    Phase = "N"
	L1N  Phase = "L1-N"
	L2N  Phase = "L2-N"
	L3N  Phase = "L3-N"
	L1L2 Phase = "L1-L2"
	L2L3 Phase = "L2-L3"
	L3L1 Phase = "L3-L1"
)

type Measurand string

const (
	energyActiveExportRegister   Measurand = "Energy.Active.Export.Register"
	energyActiveImportRegister   Measurand = "Energy.Active.Import.Register"
	energyReactiveExportRegister Measurand = "Energy.Reactive.Export.Register"
	energyReactiveImportRegister Measurand = "Energy.Reactive.Import.Register"
	energyActiveExportInterval   Measurand = "Energy.Active.Export.Interval"
	energyActiveImportInterval   Measurand = "Energy.Active.Import.Interval"
	energyReactiveExportInterval Measurand = "Energy.Reactive.Export.Interval"
	energyReactiveImportInterval Measurand = "Energy.Reactive.Import.Interval"
	powerActiveExport            Measurand = "Power.Active.Export"
	powerActiveImport            Measurand = "Power.Active.Import"
	powerOffered                 Measurand = "Power.Offered"
	powerReactiveExport          Measurand = "Power.Reactive.Export"
	powerReactiveImport          Measurand = "Power.Reactive.Import"
	powerFactor                  Measurand = "Power.Factor"
	currentImport                Measurand = "Current.Import"
	currentExport                Measurand = "Current.Export"
	currentOffered               Measurand = "Current.Offered"
	voltage                      Measurand = "Voltage"
	frequency                    Measurand = "Frequency"
	temperature                  Measurand = "Temperature"
	soc                          Measurand = "SoC"
	rpm                          Measurand = "RPM"
)

type ValueFormat string

const (
	rawFormat  ValueFormat = "Raw"
	signFormat ValueFormat = "SignedData"
)

type ReadingContext string

const (
	interruptionBegin ReadingContext = "Interruption.Begin"
	interruptionEnd   ReadingContext = "Interruption.End"
	sampleClock       ReadingContext = "Sample.Clock"
	samplePeriodic    ReadingContext = "Sample.Periodic"
	transcationBegin  ReadingContext = "Transaction.Begin"
	transcationEnd    ReadingContext = "Transaction.End"
	trigger           ReadingContext = "Trigger"
	other             ReadingContext = "Other"
)

type SampledValue struct {
	Value     string         `json:"value" validate:"required"`
	Context   ReadingContext `json:"context,omitempty" validate:"omitempty,readingContext"`
	Format    ValueFormat    `json:"format,omitempty" validate:"omitempty,valueFormat"`
	Measurand Measurand      `json:"measurand,omitempty" validate:"omitempty,measurand"`
	Phase     Phase          `json:"phase,omitempty" validate:"omitempty,phase"`
	Location  Location       `json:"location,omitempty" validate:"omitempty,location"`
	Unit      UnitOfMeasure  `json:"unit,omitempty" validate:"omitempty,unitOfMeasure"`
}

type MeterValue struct {
	TimeStamp    string         `json:"timestamp"    validate:"required,dateTime"`
	SampledValue []SampledValue `json:"sampledValue" validate:"required,min=1,dive"`
}

type MeterValuesRequest struct {
	ConnectorId   *int         `json:"connectorId" validate:"required,gte=0"`
	TransactionId *int         `json:"transactionId,omitempty" validate:"omitempty"`
	MeterValue    []MeterValue `json:"meterValue"    validate:"required,min=1,dive"`
}

func (MeterValuesRequest) Action() string {
	return MeterValuesName
}
func (r *MeterValuesRequest) Reset() {
	r.ConnectorId = nil
	r.TransactionId = nil
	r.MeterValue = nil
}

type MeterValuesResponse struct{}

func (MeterValuesResponse) Action() string {
	return MeterValuesName
}

func (r *MeterValuesResponse) Reset() {}
