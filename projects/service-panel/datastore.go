package main

import (
	"context"
	"sync"
)

func systemById(ctx context.Context, id string) *System {
	_, span := tracer.Start(ctx, "getById") // tracing
	defer span.End()                        // tracing
	if cached, ok := cache.Load(id); ok {
		if alarm, ok := cached.(*System); ok {
			return alarm
		}
	}
	for i := range database {
		if database[i].ID == id {
			cache.Store(id, &database[i])
			return &database[i]
		}
	}
	return nil
}

var cache sync.Map
var database = []System{
	{"4afe3656-f49b-4cdd-a677-c9d71f4e30b9", true, true},
	{"d227da7a-e504-4714-a100-6c7d43455466", true, true},
	{"445314cb-e37b-4337-a46b-dc65a088c953", true, true},
	{"b9200962-2585-443b-9c3c-e9ec99f8db20", true, true},
	{"4f82c2af-486c-411c-a533-a411d14fb2bb", true, true},
	{"d888a098-a548-4d44-8a26-128f8d3b66e1", true, true},
	{"167acac7-d0b7-4dcf-a55c-9ba3648b15e1", false, false},
	{"b97b9089-2d72-4916-8c91-b2fe1d7e1abf", false, false},
	{"42348e83-1125-4f00-9da4-8bb525067dd1", false, false},
	{"a45490a3-57ec-4c77-b66d-91d4c89f7785", false, false},
	{"efaee79a-2f6e-47f3-82b2-3e862c9285fb", true, false},
	{"b7f5a617-cf61-4df0-a2c4-90d9b8d01060", false, false},
	{"c51b21d5-d45c-47b5-9850-c1db21a8ef9f", false, false},
	{"1ae62da0-548b-4ba7-9bd7-1e79c030e914", false, false},
	{"ab616c7b-a67e-4f73-885f-2528b824f74f", false, false},
	{"fc98da16-3ff6-4c79-8d2d-9ad10cdd4903", false, false},
	{"3f805c89-96ba-48ef-91cc-58a18b709e99", false, false},
	{"1c04215e-24c2-43e6-bf7e-6e80c1b5a5a2", false, false},
	{"00792a24-0003-4d46-9851-4a71273cf7ad", false, false},
	{"d86b76b4-3a53-44b0-b69d-c5452e4f02d8", false, false},
	{"c20f462d-0414-4297-bbe4-1228a65b2ebd", true, false},
	{"28e85c5d-dfaa-4f02-b30a-caf9c0c04d3b", false, false},
	{"56e19633-c42a-49d3-b899-bc91e9f37419", false, false},
	{"a96c3b29-4107-4021-8e01-75a60654d8e6", true, true},
	{"b00c8c3d-fb05-4396-8b8e-f832ee90d8e5", false, false},
	{"7feadadf-7654-45e4-8765-bfd4de1893fd", false, false},
	{"d09d2d95-260d-41f6-8492-42e09e736d6f", false, false},
	{"d4770ba5-7d2a-499f-87d2-115cc09ffce2", false, false},
	{"00d27610-3969-4328-9215-42d6d55bb8c9", true, true},
	{"34dfdb7c-4d37-488f-812d-8576cf8f2730", false, false},
	{"5b3df1e4-e41c-4fb9-91de-e9bb7fbc5cda", false, false},
	{"8bb8cb8d-0ff4-4c8d-9c7f-3c9dc27c1ea6", false, false},
	{"65d14b32-f588-4d47-9f0a-9c673bbf4a63", true, false},
	{"5a0ab2c1-9f25-4e5d-b578-e9f7401e93d1", false, false},
	{"2db7a2f7-bb0e-45e7-8544-7155e9f45228", false, false},
	{"04686c28-0e36-4e7e-8b6a-f4537bb8b2f6", false, false},
	{"416c497b-202f-4452-9b9e-14a30a9b6431", true, true},
	{"3a488f54-d0c3-4881-830f-69b5c51d9ab1", false, false},
	{"3dbd3182-6d86-4454-bbcf-b5d049c2461f", false, false},
	{"f741ad5f-05b2-4b39-9b65-df9e6da50780", false, false},
	{"a1b293f8-11de-4e07-b921-c9379c8b5a2f", true, true},
	{"2f839e33-4711-46b2-83db-8f6f6a46f1e1", false, false},
}
