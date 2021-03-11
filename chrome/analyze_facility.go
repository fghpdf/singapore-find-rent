package chrome

func analyzeFac(facilities []string) *Facility {
	fac := &Facility{}
	for _, facility := range facilities {
		switch facility {
		case "WADING POOL", "50M LAP POOL", "JACUZZI POOL", "50M FREEFORM POOL", "BEACH SPLASH POOL", "FAMILY POOL", "REFLECTION POOL":
			fac.Pool = true
		case "TENNIS COURT":
			fac.TennisCourt = true
		case "READING Corner":
			fac.ReadingCorner = true
		case "FITNESS STATION", "FITNESS ALCOVE":
			fac.FitnessArea = true
		case "INDOOR GYM", "HYDRO GYM STATION":
			fac.Gymnasium = true
		case "BBQ AREA":
			fac.BbqPit = true
		case "24-HOUR SECURITY":
			fac.Security = true
		}
	}

	return fac
}
