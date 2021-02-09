package amdfw

import (
	"bytes"
	"fmt"
	"log"
)

const DefaultFlashMapping1 = uint32(0xFF000000)
const DefaultFlashMapping2 = uint32(0x00000000)

func GetFlashMapping(firmwareBytes []byte, fet *FirmwareEntryTable) (uint32, error) {

	type mappingMagic struct {
		addr  *uint32
		magic []string
	}
	var err error
	var mapping uint32

	for _, s := range []mappingMagic{{
		addr:  fet.PSPDirBase,
		magic: []string{PSPCOOCKIE, DUALPSPCOOCKIE},
	}, {
		addr:  fet.NewPSPDirBase,
		magic: []string{PSPCOOCKIE, DUALPSPCOOCKIE},
	}, {
		addr:  fet.BHDDirBase,
		magic: []string{BHDCOOCKIE},
	}, {
		addr:  fet.BHDDirBase,
		magic: []string{BHDCOOCKIE},
	},
	} {
		for _, m := range s.magic {

			if s.addr != nil && *s.addr != 0 {
				mapping, err = testMapping(firmwareBytes, *s.addr, m)
				if err == nil {
					return mapping, nil
				}
			}
		}
	}
	return 0, fmt.Errorf("No valid mapping found: %v", err)
}

func testMapping(firmwareBytes []byte, address uint32, expected string) (uint32, error) {

	for _, mapping := range []uint32{
		DefaultFlashMapping1 + 0x000000, //16M
		DefaultFlashMapping1 + 0x800000, // 8M
		DefaultFlashMapping1 + 0xB00000, // 4M
		DefaultFlashMapping1 + 0xD00000, // 2M
		DefaultFlashMapping1 + 0xE00000, // 1M
		DefaultFlashMapping1 + 0xE80000, // 512K
		DefaultFlashMapping2 + 0x000000, //16M
		DefaultFlashMapping2 + 0x800000, // 8M
		DefaultFlashMapping2 + 0xB00000, // 4M
		DefaultFlashMapping2 + 0xD00000, // 2M
		DefaultFlashMapping2 + 0xE00000, // 1M
		DefaultFlashMapping2 + 0xE80000, // 512K
	} {

		expectedBytes := []byte(expected)
		testAddr := address - mapping
		if int(testAddr) > len(firmwareBytes) {
			continue
		}

		log.Printf(
			"trying: %v, %v - %v",
			testAddr,
			address,
			mapping,
		)
		if bytes.Equal(firmwareBytes[testAddr:testAddr+4], expectedBytes) {
			return mapping, nil
		}
		log.Printf(
			"not matching: %v is not %v",
			firmwareBytes[testAddr:testAddr+4],
			expectedBytes,
		)
	}
	return 0, fmt.Errorf("No Default Mapping fits")
}
