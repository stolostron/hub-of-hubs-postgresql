package compliance

import (
	"github.com/gofrs/uuid"
)

const policiesNumber = 100

//nolint:gochecknoglobals
var policyUUIDs []uuid.UUID

func mustFromString(s string) uuid.UUID {
	uuid, err := uuid.FromString(s)
	if err != nil {
		panic("wrong format for the UUID")
	}

	return uuid
}

//nolint: gochecknoinits, funlen
func init() {
	policyUUIDs = []uuid.UUID{
		mustFromString("2fe63eca-ee41-11eb-b83c-37beaee817b2"),
		mustFromString("2fe6fbee-ee41-11eb-92e3-db8718063d08"),
		mustFromString("2fe7b282-ee41-11eb-ba1e-bb934b066872"),
		mustFromString("2fe86b5a-ee41-11eb-8047-172c7aff65eb"),
		mustFromString("2fe92630-ee41-11eb-bc56-9f0e85d312c3"),
		mustFromString("2fe9dc10-ee41-11eb-9df7-a3e75727b174"),
		mustFromString("2fea940c-ee41-11eb-87ab-931256dc29a7"),
		mustFromString("2feb578e-ee41-11eb-b85b-8761d5aec20f"),
		mustFromString("2fec135e-ee41-11eb-baf7-0b945cc364a0"),
		mustFromString("2feccb96-ee41-11eb-b571-73579cacbc46"),
		mustFromString("2fed86da-ee41-11eb-9ec5-23d6b0ead578"),
		mustFromString("2fee3ea4-ee41-11eb-a71c-ffebc3851341"),
		mustFromString("2feef560-ee41-11eb-94c1-af642ce0bb45"),
		mustFromString("2fefad98-ee41-11eb-b149-33fd90fe88f8"),
		mustFromString("2ff069c2-ee41-11eb-bd30-573dff337508"),
		mustFromString("2ff12498-ee41-11eb-8a61-f772850b03d0"),
		mustFromString("2ff1df1e-ee41-11eb-851d-3b4b4f7354bb"),
		mustFromString("2ff2a25a-ee41-11eb-9092-0b269499407c"),
		mustFromString("2ff36596-ee41-11eb-b766-6bbe62df9c86"),
		mustFromString("2ff42832-ee41-11eb-8aad-b72570dfe83d"),
		mustFromString("2ff4eb5a-ee41-11eb-8f7a-f376eb5de35d"),
		mustFromString("2ff5af86-ee41-11eb-a526-171c852f0eae"),
		mustFromString("2ff66ffc-ee41-11eb-be6e-331668c68bb8"),
		mustFromString("2ff72910-ee41-11eb-8064-afd5d0475390"),
		mustFromString("2ff7f6ce-ee41-11eb-8e4d-ef024db968b9"),
		mustFromString("2ff8b578-ee41-11eb-9f65-8f8d03d9091a"),
		mustFromString("2ff96edc-ee41-11eb-b08b-eb8e401f08cf"),
		mustFromString("2ffa3178-ee41-11eb-8501-6319d548e396"),
		mustFromString("2ffaf1f8-ee41-11eb-a265-57465af36abd"),
		mustFromString("2ffbb98a-ee41-11eb-aaa6-93b61cec9db0"),
		mustFromString("2ffc7834-ee41-11eb-aeed-7b74c95aca9b"),
		mustFromString("2ffd3cce-ee41-11eb-a66e-d3255db9b953"),
		mustFromString("2ffdf5b0-ee41-11eb-9919-b7e0375d055f"),
		mustFromString("2ffeb81a-ee41-11eb-8e55-6f28b8763c6a"),
		mustFromString("2fff7e1c-ee41-11eb-8bb4-1b7d22a0787d"),
		mustFromString("30003988-ee41-11eb-82bf-1fed0bc46683"),
		mustFromString("3000f29c-ee41-11eb-83a8-b7d4c7dd7bce"),
		mustFromString("3001b2c2-ee41-11eb-86a2-a31db5988934"),
		mustFromString("30027d74-ee41-11eb-a607-43007be55441"),
		mustFromString("30033e62-ee41-11eb-80cb-bb5f864f8aef"),
		mustFromString("30042796-ee41-11eb-9de2-cfe88b5b5d1e"),
		mustFromString("300504cc-ee41-11eb-8788-333bd16ed795"),
		mustFromString("3005d5d2-ee41-11eb-bf1b-13bd242fd792"),
		mustFromString("300696c0-ee41-11eb-90cb-d3819e3a1250"),
		mustFromString("30075b8c-ee41-11eb-af61-6bd50f1b1989"),
		mustFromString("30081554-ee41-11eb-a63d-8f13cc7f33e5"),
		mustFromString("3008cf76-ee41-11eb-bd65-af862938bf6a"),
		mustFromString("30098efc-ee41-11eb-8c8d-03e668337ae8"),
		mustFromString("300a4a90-ee41-11eb-ae87-9798100947ee"),
		mustFromString("300b0502-ee41-11eb-ae9b-fb33cd4f81aa"),
		mustFromString("300bd180-ee41-11eb-a190-0fb1bea6f58e"),
		mustFromString("300ce228-ee41-11eb-866d-37a7e1f6c195"),
		mustFromString("300de100-ee41-11eb-b002-5b29b215a343"),
		mustFromString("300ec8a4-ee41-11eb-972d-33ee60ba9dc7"),
		mustFromString("300fd6d6-ee41-11eb-a8a7-ff23993b0d3e"),
		mustFromString("3010c816-ee41-11eb-89e2-bfacd4efd567"),
		mustFromString("3011be92-ee41-11eb-ade6-97ecd49ac080"),
		mustFromString("3012be50-ee41-11eb-be4b-03a0bf132f84"),
		mustFromString("3013d4c0-ee41-11eb-9eb7-e3e605d3787d"),
		mustFromString("301545da-ee41-11eb-be18-6f0186bee8c9"),
		mustFromString("301610b4-ee41-11eb-b915-63d88cba115b"),
		mustFromString("30172c38-ee41-11eb-be12-8fd6c64d3e9d"),
		mustFromString("3018172e-ee41-11eb-86c6-17139cfcadb9"),
		mustFromString("301908a0-ee41-11eb-bae0-2b9020823e52"),
		mustFromString("301a19de-ee41-11eb-b7f3-831ba11cefcb"),
		mustFromString("301b051a-ee41-11eb-a1f7-2320fcc8a72e"),
		mustFromString("301bdbac-ee41-11eb-b7aa-a71599e5345e"),
		mustFromString("301ca0d2-ee41-11eb-98a6-af13a889bd83"),
		mustFromString("301d5f7c-ee41-11eb-91ec-ff8f705f3037"),
		mustFromString("301e2aa6-ee41-11eb-b3f8-5706c045e055"),
		mustFromString("301eee96-ee41-11eb-9177-5f1ab7411e95"),
		mustFromString("301faeb2-ee41-11eb-97d5-f3b6f92a3f71"),
		mustFromString("30206d5c-ee41-11eb-9970-3310ad901280"),
		mustFromString("302131ce-ee41-11eb-a4bc-fbfe489acb03"),
		mustFromString("3021fbcc-ee41-11eb-81e2-930a224f89e7"),
		mustFromString("3022d916-ee41-11eb-9006-9b7bb5d7d0a1"),
		mustFromString("3023a0bc-ee41-11eb-88ed-0f62280d56bd"),
		mustFromString("30246218-ee41-11eb-b93c-f35dd853fb97"),
		mustFromString("302529b4-ee41-11eb-a7e8-ff45ca558275"),
		mustFromString("3025f1c8-ee41-11eb-86df-dfd5d93aa1ed"),
		mustFromString("3026af14-ee41-11eb-947b-cb65f63ab07f"),
		mustFromString("30276b84-ee41-11eb-8d42-137edd83a9b3"),
		mustFromString("30282d76-ee41-11eb-8a58-97b17fcba0cb"),
		mustFromString("3028f33c-ee41-11eb-9e08-5f9f460f0772"),
		mustFromString("3029b33a-ee41-11eb-971e-7b0352f66f4a"),
		mustFromString("302a75f4-ee41-11eb-bc21-efc940026d20"),
		mustFromString("302b353e-ee41-11eb-895c-2fae92366322"),
		mustFromString("302bef7e-ee41-11eb-9fd6-d3a648869049"),
		mustFromString("302cb49a-ee41-11eb-9857-c35546e6ca02"),
		mustFromString("302d7466-ee41-11eb-9988-672a6651659c"),
		mustFromString("302e36c6-ee41-11eb-864f-771fc95ffa4f"),
		mustFromString("302ef8e0-ee41-11eb-afe8-df0f5f370c04"),
		mustFromString("302fbde8-ee41-11eb-a304-d3ea3683ecf4"),
		mustFromString("303080ac-ee41-11eb-ba3b-1f0ddeada85e"),
		mustFromString("30313de4-ee41-11eb-b668-d31ada7067b4"),
		mustFromString("303203b4-ee41-11eb-bee6-2fb7e80bdbcb"),
		mustFromString("3032d014-ee41-11eb-8ad5-ffa736973c4a"),
		mustFromString("303397e2-ee41-11eb-b36f-0bf1fece9e8b"),
		mustFromString("30348e04-ee41-11eb-8a76-0792792772f3"),
		mustFromString("3035631a-ee41-11eb-b86a-9761c82dc6a7"),
	}
}
