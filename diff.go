package genv

type ChangeValue struct {
	NewValue string
	OldValue string
}

type Diff struct {
	Added   map[string]string
	Removed map[string]string
	Changed map[string]ChangeValue
}

func DiffEnvMap(old, new map[string]string) Diff {
	added := make(map[string]string)
	removed := make(map[string]string)
	changed := make(map[string]ChangeValue)

	for k, v := range old {
		if _, ok := new[k]; !ok {
			removed[k] = v
		} else if new[k] != v {
			changed[k] = ChangeValue{NewValue: new[k], OldValue: v}
		}
	}

	for k, v := range new {
		if _, ok := old[k]; !ok {
			added[k] = v
		}
	}

	return Diff{
		Added:   added,
		Removed: removed,
		Changed: changed,
	}
}
