package bconfig

// deepMerge recursively merges src into dst.
// - If both values are maps, merge recursively.
// - If source value is scalar, override destination value.
// - If source value is slice, override destination slice.
// - If source value is nil, override destination value with nil.
func deepMerge(dst, src map[string]any) map[string]any {
	if dst == nil {
		dst = make(map[string]any)
	}

	for k, srcVal := range src {
		if srcVal == nil {
			dst[k] = nil
			continue
		}

		dstVal, exists := dst[k]
		if !exists {
			if srcMap, ok := toMap(srcVal); ok {
				dst[k] = deepMerge(nil, srcMap)
			} else {
				dst[k] = cloneValue(srcVal)
			}
			continue
		}

		srcMap, srcIsMap := toMap(srcVal)
		dstMap, dstIsMap := toMap(dstVal)

		if srcIsMap && dstIsMap {
			dst[k] = deepMerge(dstMap, srcMap)
		} else {
			// Either one of them is not a map, or neither.
			// Override destination value with source value.
			if srcIsMap {
				dst[k] = deepMerge(nil, srcMap)
			} else {
				dst[k] = cloneValue(srcVal)
			}
		}
	}

	return dst
}

// toMap converts a value to map[string]any if possible.
func toMap(v any) (map[string]any, bool) {
	if m, ok := v.(map[string]any); ok {
		return m, true
	}
	return nil, false
}

// cloneValue performs a deep copy of maps and slices to ensure the loaded config
// is independent of the source structures.
func cloneValue(v any) any {
	if v == nil {
		return nil
	}

	if m, ok := toMap(v); ok {
		return deepMerge(nil, m)
	}

	if slice, ok := v.([]any); ok {
		res := make([]any, len(slice))
		for i, val := range slice {
			res[i] = cloneValue(val)
		}
		return res
	}

	// For basic types, they are copied by value, so return directly.
	return v
}
