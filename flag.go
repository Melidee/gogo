package main

type Flag struct {
	Name       string
	Short      rune
	Long       string
	About      string
	TakesValue bool
	HasDefault bool
	Default    string
	Required   bool
}

func NewFlag(name string) *Flag {
	return &Flag{
		Name: name,
	}
}

func (f *Flag) SetShort(short rune) *Flag {
	f.Short = short
	return f
}

func (f *Flag) SetLong(long string) *Flag {
	f.Long = long
	return f
}

func (f *Flag) SetAbout(about string) *Flag {
	f.About = about
	return f
}

func (f *Flag) SetDefault(defaultValue string) *Flag {
	f.HasDefault = true
	f.Default = defaultValue
	return f
}

func (f *Flag) SetRequired(yes bool) *Flag {
	f.Required = yes
	return f
}

func (f *Flag) Matches(iter *ArgIterator) *FlagMatch {
	if !iter.HasNext() {
		return nil
	}
	flag := iter.Peek()
	if len(flag) < 2 {
		return nil
	}
	flag_is_long := flag[0:2] == "--" && flag[2:] == f.Long
	flag_is_short := len(flag) == 2 && flag[0] == '-' && flag[1:] == string(f.Short)
	if !(flag_is_long || flag_is_short) {
		return nil
	}
	iter.Next() // consume flag
	match := FlagMatch{
		Name: f.Name,
	}
	if !f.TakesValue {
		return &match		
	}
	if !iter.HasNext() {
		if f.HasDefault {
			
		}
	}
	return &match
}

type FlagMatch struct {
	Name     string
	HasValue bool
	Value    string
}
