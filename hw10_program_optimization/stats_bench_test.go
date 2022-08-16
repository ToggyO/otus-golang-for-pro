package hw10programoptimization

import (
	"archive/zip"
	"flag"
	"testing"
)

var isOld = flag.Bool("old", false, "bench non optiomized version of GetDomainStat (default is false)")

func BenchmarkGetDomainStat(b *testing.B) {
	domain := "biz"
	becnhingFuncName := "GetDomainStat"
	if *isOld {
		becnhingFuncName = "GetDomainStatOld"
	}

	b.Logf("benching %s", becnhingFuncName)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		reader, err := zip.OpenReader("./testdata/users.dat.zip")
		if err != nil {
			b.Fatal()
		}

		file, err := reader.File[0].Open()
		if err != nil {
			b.Fatal()
		}

		b.StartTimer()
		if *isOld {
			_, _ = GetDomainStatOld(file, domain)
			continue
		}

		_, _ = GetDomainStat(file, domain)
		_ = reader.Close()
	}
}
