package models

type BankBranch struct {
	Name              string `json:"name"`
	FK_Bank_Id        string `json:"bank_Id"`
	Branch            string `json:"branch"`
	Phone             uint   `json:"phone"`
	Email             string `json:"email" `
	PK_Bank_Branch_Id string `json:"bank_Branch_Id" gorm:"primaryKey"`
	Address           string `json:"address"`
	Bank_IFSC         string `json:"bank_IFSC"`
}

// //SAMPLE REQUEST
// {
//   "name": "Unity",
// "branch": "Uttam Nagar",
// "bank_Id":"UNIT0002",
// "phone":12345654,
// "email":"unityuttam@gmail.com",
// "address":"Najafgarh Road, Metro Pillar: 687, ND-110059"
// }

// //SAMPLE RESPONSE
// {
// "name":"Unity",
// "bank_Id":"UNIT0002",
// "branch":"Uttam Nagar",
// "phone":12345654,
// "email":"unityuttam@gmail.com",
// "bank_Branch_Id":"UTUTTAM NAGAR020001",
// "address":"Najafgarh Road, Metro Pillar: 687, ND-110059"
// "bank_IFSC":"UNIT0000001"
//}
