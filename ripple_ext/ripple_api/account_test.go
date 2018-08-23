package ripple_api

import "testing"

func TestAccount(t *testing.T)  {
	/// test0
	ss1 := make(map[string]*RippledKeyPair)
	for i := 0; i < 1000000 ;i++ {
		k, err := NewAccount()
		if err != nil {
			t.Fatalf("NewAccount err : %v", err)
		}

		if _, ok := ss1[k.Seed]; ok {
			t.Fatalf("%s is exist", k.Seed)
		}

		ss1[k.Seed] = k
	}
	for _, k := range ss1 {
		addr, err := GetAddressBySeed(k.Seed)
		if err != nil {
			t.Fatalf("%s GetAddressBySeed err : %v", k.Seed, err)
		}

		if addr != k.Address {
			t.Fatalf("%s is not the seed of %s", k.Seed, k.Address)
		}

		t.Log(k.Seed, "--", k.Address)
	}

	/// test1
	ss := []RippledKeyPair{
		RippledKeyPair{
			Address:"rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh",
			Seed:"snoPBrXtMeMyMHUVTgbuqAfg1SUTb",
		},
		RippledKeyPair{
			Address:"rEiVPf31YhtgJtq8fvi1cpwNKDEaAp25Ud",
			Seed:"ssLUvvnEDgEjGTJorSB5aMscebBbt",
		},
		RippledKeyPair{
			Address:"rBCSgao1iBP412bpFrQgfVnj7HGEeJUThr",
			Seed:"sh6JW3F9EtJ7iA8ZHLg85bcDVqFhy",
		},
		RippledKeyPair{
			Address:"rntihVJZdoWTchfqgKRFGBK2byUUQ6iXze",
			Seed:"shEpDkHDQpPFPWsmgzCLMzaBWVF92",
		},
	}

	for _, k := range ss {
		addr, err := GetAddressBySeed(k.Seed)
		if err != nil {
			t.Fatalf("%s GetAddressBySeed err : %v", k.Seed, err)
		}

		if addr != k.Address {
			t.Fatalf("%s is not the seed of %s", k.Seed, k.Address)
		}

		t.Log(k.Seed, "--", k.Address, "--", addr)
	}
}
