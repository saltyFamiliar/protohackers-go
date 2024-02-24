package means2anend

import (
	"fmt"
	"log"
	"time"
)

type PriceHistory []RequestPacket

func InsertRequestPacket(ph PriceHistory, rp RequestPacket) (PriceHistory, error) {
	newPriceHistory := make(PriceHistory, len(ph)+1)

	if len(ph) == 0 {
		newPriceHistory[0] = rp
		return newPriceHistory, nil
	}

	timeStamp := rp.TimeStamp()
	var insertionPoint int
	for ; insertionPoint < len(ph); insertionPoint++ {
		if ph[insertionPoint].TimeStamp() > timeStamp {
			break
		} else if ph[insertionPoint].TimeStamp() == timeStamp {
			return nil, fmt.Errorf("duplicate timestamp found")
		}
		newPriceHistory[insertionPoint] = ph[insertionPoint]
	}

	for i := insertionPoint; i < len(ph); i++ {
		if ph[i].TimeStamp() == timeStamp {
			return nil, fmt.Errorf("duplicate timestamp found")
		}
		newPriceHistory[i+1] = ph[i]
	}

	newPriceHistory[insertionPoint] = rp

	return newPriceHistory, nil
}

func (ph PriceHistory) ShowAll() {
	for _, rp := range ph {
		timestamp := time.Unix(int64(rp.TimeStamp()), 0)
		log.Printf("Timestamp: %v, Price: %d", timestamp, rp.Price())
	}
}

func (ph PriceHistory) InRange(start, end int32) (pInRange PriceHistory) {
	for _, rp := range ph {
		timestamp := rp.TimeStamp()
		if timestamp >= start && timestamp <= end {
			pInRange = append(pInRange, rp)
		}
	}
	return pInRange
}

func (ph PriceHistory) MeanPrice() int32 {
	if len(ph) == 0 {
		return 0
	}

	var total int64
	for _, rp := range ph {
		total += int64(rp.Price())
	}

	return int32(total / int64(len(ph)))
}
