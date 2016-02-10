package mt

type MeasWriter interface {
	New(size int64)
	Add(m Meas) bool
	Add_range(m []Meas) int64
	Cap() int64
	IsFull() bool
	Close()
}

type MeasReader interface {
	ReadAll() []Meas
	Read(ids []Id, from, to Time) []Meas
	ReadFltr(ids []Id, flg Flag, from, to Time) []Meas
	TimePoint(ids []Id, time Time) []Meas
	TimePointFltr(ids []Id, flg Flag, time Time) []Meas
}

type MeasStorage interface {
	MeasReader
	MeasWriter
}

type DataReader interface {
	Read(m *[]Meas, count int64)
	ReadAll() []Meas
}
