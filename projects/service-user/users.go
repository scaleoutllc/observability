package main

import (
	"context"
	"sync"
)

type User struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Email   string   `json:"email"`
	Country string   `json:"country"`
	Systems []string `json:"systems"`
}

func getAllUsers(ctx context.Context) []User {
	_, span := tracer.Start(ctx, "getAllUsers") // tracing
	defer span.End()                            // tracing
	return database
}

func getUserByID(ctx context.Context, id string) *User {
	_, span := tracer.Start(ctx, "getUserByID") // tracing
	defer span.End()                            // tracing
	if user, ok := cache.Load(id); ok {
		return user.(*User)
	}
	for _, user := range database {
		if user.ID == id {
			cache.Store(id, &user)
			return &user
		}
	}
	return nil
}

func getUserBySystemID(ctx context.Context, id string) *User {
	_, span := tracer.Start(ctx, "getUserBySystemID") // tracing
	defer span.End()                                  // tracing
	for _, user := range database {
		for _, panel := range user.Systems {
			if panel == id {
				return &user
			}
		}
	}
	return nil
}

var cache sync.Map
var database = []User{
	{"c0e8e5da-cd7f-4ed5-9a1f-df27a2d24592", "Tomás García", "tomas.garcia@correo.es", "Spain", []string{"4afe3656-f49b-4cdd-a677-c9d71f4e30b9"}},
	{"841e22ce-39c3-4882-8e50-6d74c05c728d", "María Sánchez", "maria.sanchez@outlook.es", "Spain", []string{"d227da7a-e504-4714-a100-6c7d43455466"}},
	{"2acee4b0-241c-4f7b-8ed3-8cdee5276586", "Alejandro García", "alejandro.garcia@ymail.es", "Spain", []string{"445314cb-e37b-4337-a46b-dc65a088c953"}},
	{"28149ad2-2a9e-4888-a3e5-ce470e7625a5", "Lidia Fernández", "lidia.fernandez@live.es", "Spain", []string{"b9200962-2585-443b-9c3c-e9ec99f8db20"}},
	{"7acadf12-9345-4c35-8cc7-6abcc8bea31e", "David López", "david.lopez@rocketmail.es", "Spain", []string{"4f82c2af-486c-411c-a533-a411d14fb2bb"}},
	{"da39aafd-20a5-45ee-9f44-6b7972fc9a23", "Henrik Andersson", "henrik.andersson@live.se", "Sweden", []string{"d888a098-a548-4d44-8a26-128f8d3b66e1"}},
	{"bf3a7a62-3350-4749-88cb-0ef5de2ba17e", "Chloé Lefevre", "chloe.lefevre@wanadoo.fr", "France", []string{"167acac7-d0b7-4dcf-a55c-9ba3648b15e1", "b97b9089-2d72-4916-8c91-b2fe1d7e1abf"}},
	{"846d3e01-2674-49f0-ae7b-f58a4f1f76ec", "Jasmine Lindberg", "jasmine.lindberg@yahoo.se", "Sweden", []string{"42348e83-1125-4f00-9da4-8bb525067dd1"}},
	{"d122efb2-c01d-4227-8f10-952baa388a91", "Emma Johansson", "emma.johansson@rocketmail.se", "Sweden", []string{"a45490a3-57ec-4c77-b66d-91d4c89f7785"}},
	{"c06a9758-2a8a-4276-ac5d-897a8be63737", "Vittoria Rossi", "vittoria.rossi@yahoo.it", "Italy", []string{"efaee79a-2f6e-47f3-82b2-3e862c9285fb"}},
	{"cd4d0a2b-97c9-4a2a-bcc7-aa5d1beb6541", "Michela Bianchi", "michela.bianchi@hotmail.it", "Italy", []string{"b7f5a617-cf61-4df0-a2c4-90d9b8d01060"}},
	{"a159fc07-061f-48a9-91d9-529132f81e70", "Tyler Bellini", "tyler.bellini@yahoo.it", "Italy", []string{"c51b21d5-d45c-47b5-9850-c1db21a8ef9f"}},
	{"6710a6fc-d752-45f4-a526-a6351c5f8a64", "Isabella O'Connor", "isabella.oconnor@ymail.ie", "Ireland", []string{"1ae62da0-548b-4ba7-9bd7-1e79c030e914"}},
	{"7c34a32c-5c94-494f-8a9b-2b9b9f646d6d", "Dylan O'Sullivan", "dylan.osullivan@outlook.ie", "Ireland", []string{"ab616c7b-a67e-4f73-885f-2528b824f74f"}},
	{"93c68315-376a-4df2-915b-9a45030b911e", "Anna Alexander", "anna.alexander@yahoo.ie", "Ireland", []string{"fc98da16-3ff6-4c79-8d2d-9ad10cdd4903"}},
	{"5c909e5d-b706-48db-ae8b-c3ba5c216e89", "Aisling Kelly", "aisling.kelly@outlook.ie", "Ireland", []string{"3f805c89-96ba-48ef-91cc-58a18b709e99", "1c04215e-24c2-43e6-bf7e-6e80c1b5a5a2"}},
	{"d90bcd4b-37e9-4858-9bfb-b81c979adccd", "Hugo Dubois", "hugo.dubois@outlook.fr", "France", []string{"00792a24-0003-4d46-9851-4a71273cf7ad", "d86b76b4-3a53-44b0-b69d-c5452e4f02d8"}},
	{"b84e6a59-3b40-4df9-83e0-3247cbff25b1", "Daniel O'Connell", "daniel.oconnell@rocketmail.ie", "Ireland", []string{"c20f462d-0414-4297-bbe4-1228a65b2ebd", "28e85c5d-dfaa-4f02-b30a-caf9c0c04d3b"}},
	{"57efbff5-c520-4d5a-ae30-53b73c6d7f2a", "Brianna O'Neill", "brianna.oneill@rocketmail.ie", "Ireland", []string{"56e19633-c42a-49d3-b899-bc91e9f37419"}},
	{"fd1dc982-49de-4a20-98f1-bd1a6cf0375f", "Théo Dupont", "theo.dupont@gmail.fr", "France", []string{"a96c3b29-4107-4021-8e01-75a60654d8e6"}},
	{"c29f83e5-5db4-41f7-8a18-1a1d9f5018c5", "Julien Moreau", "julien.moreau@aol.fr", "France", []string{"b00c8c3d-fb05-4396-8b8e-f832ee90d8e5"}},
	{"ff748e7c-621e-4c6b-905f-99d7c6d2c82f", "Ashley Dubois", "ashley.dubois@hotmail.fr", "France", []string{"7feadadf-7654-45e4-8765-bfd4de1893fd"}},
	{"68299257-4bc2-4041-b43e-f66b29d1d379", "Léa Martin", "lea.martin@wanadoo.fr", "France", []string{"d09d2d95-260d-41f6-8492-42e09e736d6f", "d4770ba5-7d2a-499f-87d2-115cc09ffce2"}},
	{"ed5e883a-7c47-42df-9334-6c53f1b61005", "Richard Dubois", "richard.dubois@aol.fr", "France", []string{"00d27610-3969-4328-9215-42d6d55bb8c9"}},
	{"78e9f45b-5402-41ac-9b23-ec74b9ad5e3f", "Erik Hansen", "erik.hansen@gmail.no", "Norway", []string{"34dfdb7c-4d37-488f-812d-8576cf8f2730"}},
	{"f14c33a8-ecf2-4da0-93b1-739453e7e730", "Luca Colombo", "luca.colombo@libero.it", "Italy", []string{"5b3df1e4-e41c-4fb9-91de-e9bb7fbc5cda"}},
	{"19d2e7da-412b-4cdd-85f5-bca77b9a4e09", "Matteo Ricci", "matteo.ricci@yahoo.it", "Italy", []string{"8bb8cb8d-0ff4-4c8d-9c7f-3c9dc27c1ea6"}},
	{"fc213c47-5a0c-43c9-9267-9914b4e36a6f", "Daniela Rossi", "daniela.rossi@live.it", "Italy", []string{"65d14b32-f588-4d47-9f0a-9c673bbf4a63", "5a0ab2c1-9f25-4e5d-b578-e9f7401e93d1"}},
	{"a40d1dbf-942a-4d1a-a89c-3b5585de659a", "Jean-Luc Fournier", "jeanluc.fournier@gmail.fr", "France", []string{"2db7a2f7-bb0e-45e7-8544-7155e9f45228"}},
	{"45f8d66e-d002-4c3c-9bcf-142a94656d14", "Anders Larsen", "anders.larsen@gmail.no", "Norway", []string{"04686c28-0e36-4e7e-8b6a-f4537bb8b2f6"}},
	{"8b12f5a3-dbb4-41be-8708-c44d02f31991", "Ingrid Nilsen", "ingrid.nilsen@outlook.no", "Norway", []string{"416c497b-202f-4452-9b9e-14a30a9b6431"}},
	{"c0b8ed5f-c805-4507-bb2b-0f22c196ee24", "Elias Hansen", "elias.hansen@rocketmail.no", "Norway", []string{"3a488f54-d0c3-4881-830f-69b5c51d9ab1"}},
	{"7a290d73-f5e5-48f4-a941-c6d2a9b374d7", "Karin Berg", "karin.berg@hotmail.se", "Sweden", []string{"3dbd3182-6d86-4454-bbcf-b5d049c2461f"}},
	{"d2b83556-8327-42f3-b021-77e9f148b403", "Liam O'Brien", "liam.obrien@gmail.ie", "Ireland", []string{"f741ad5f-05b2-4b39-9b65-df9e6da50780"}},
	{"c3be4ae5-b159-4510-bdc3-f92827d26c18", "Eoghan Murphy", "eoghan.murphy@live.ie", "Ireland", []string{"a1b293f8-11de-4e07-b921-c9379c8b5a2f"}},
	{"d45a1a9f-1e1b-4483-9e37-3247a7684972", "Louis Martin", "louis.martin@gmail.fr", "France", []string{"2f839e33-4711-46b2-83db-8f6f6a46f1e1"}},
}
