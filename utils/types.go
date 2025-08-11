package utils

func IntPtrFromInt32Ptr(p *int32) *int {
	if p == nil {
		return nil
	}
	v := int(*p)
	return &v
}
