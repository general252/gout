package ulog2

type Tag []string

func (c *Tag) String() string {
	if c == nil {
		return ""
	}
	var ll = []string(*c)
	var count = len(ll)
	if count > 0 {
		var r = "["
		for i := 0; i < count-1; i++ {
			r += ll[i] + ", "
		}
		r += ll[count-1] + "] "

		return r
	}

	return ""
}

func (c *Tag) Append(tags ...string) {
	if c == nil {
		return
	}
	*c = append(*c, tags...)
}

func (c *Tag) AddAndNew(tags ...string) Tag {
	if c == nil {
		return nil
	}

	var t Tag
	t = append(t, *c...)
	t = append(t, tags...)

	return t
}

func Tags(tags ...string) Tag {
	return tags
}
