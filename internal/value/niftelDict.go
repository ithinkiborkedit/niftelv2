package value

type DictEntry struct {
	Key   Value
	Value Value
}

type NiftelDict struct {
	buckets map[uint64][]DictEntry
}

func NewNiftelDict() *NiftelDict {
	return &NiftelDict{buckets: make(map[uint64][]DictEntry)}
}

func (d *NiftelDict) Set(key Value, val Value) {
	hash := key.Hash()
	bucket := d.buckets[hash]
	for i, entry := range bucket {
		if entry.Key.Equals(key) {
			bucket[i].Value = val
			d.buckets[hash] = bucket
			return
		}
	}

	d.buckets[hash] = append(d.buckets[hash], DictEntry{Key: key, Value: val})
}

func (d *NiftelDict) Get(key Value) (Value, bool) {
	hash := key.Hash()
	bucket := d.buckets[hash]
	for _, entry := range bucket {
		if entry.Key.Equals(key) {
			return entry.Value, true
		}
	}
	return Null(), false
}

func (d *NiftelDict) Delete(key Value) bool {
	hash := key.Hash()
	bucket := d.buckets[hash]
	for i, entry := range bucket {
		if entry.Key.Equals(key) {
			d.buckets[hash] = append(bucket[:i], bucket[i+1:]...)
			return true
		}
	}
	return false
}

func (d *NiftelDict) Keys() []Value {
	keys := []Value{}
	for _, bucket := range d.buckets {
		for _, entry := range bucket {
			keys = append(keys, entry.Key)
		}
	}
	return keys
}

func (d *NiftelDict) Iter() []DictEntry {
	entries := []DictEntry{}
	for _, bucket := range d.buckets {
		entries = append(entries, bucket...)
	}
	return entries
}
