package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)
var asset_ID int=1
//
type managingJongchin struct {
	contractapi.Contract
}

type Assetstruct struct {
	Name string `json:NAME`
	Asset_ID int `json:ASSET_ID`
	Asset_type int `json:ASSET_TYPE`
	Category string `json:CATEGORY`
	Owner []AssetOwner `json:OWNER`
	//share struct
}

type AssetOwner struct {
	OwnerName string `json:owner`
	Stake int `json:int`
}
//type sharestruct struct{}
func (s *managingJongchin) AssetExists(ctx contractapi.TransactionContextInterface, Asset_ID int)(bool,error){
	// 원장에 AssetID를 키로 갖는 데이터가 존재하는지 
	AssetID := strconv.Itoa(Asset_ID)
	assetbytes, err := ctx.GetStub().GetState(AssetID)
	
	if err != nil {
		return false,err 
	}
	return assetbytes !=nil, nil  
}

// func (s *managingJongchin) FirstAsset(ctx contractapi.TransactionContextInterface, Asset_ID int)(bool,error){ 

// }

func (s *managingJongchin) AddAsset(ctx contractapi.TransactionContextInterface, name string, asset_type int, category string, owners string, share int) error {  // 피어에서 호출, BN
	
	fmt.Println("addstart")
	asset_ID+=1
	exists, err := s.AssetExists(ctx,asset_ID)
	if err != nil {
		return fmt.Errorf(`AssetExists Error`)
	}
	if exists !=false {
		return fmt.Errorf(`AssetID (%d) is already exists`,asset_ID)
	} 
	var owner = AssetOwner{OwnerName:owners ,Stake:share}
	asset := Assetstruct {}
	asset.Name = name 
	asset.Asset_ID = asset_ID
	asset.Asset_type=asset_type
	asset.Category = category
	asset.Owner= append(asset.Owner,owner)

	assetID := strconv.Itoa(asset_ID)
	assetbytes, err := json.Marshal(asset)
	if err != nil {
		return fmt.Errorf(`Marshal error`)
	}
	fmt.Println("addEnd")
	return ctx.GetStub().PutState(assetID, assetbytes)
}
// 체인코드  type managingJongchin 부터 

func (s *managingJongchin) ReadAsset(ctx contractapi.TransactionContextInterface, Asset_ID int ) (*Assetstruct,error){
	asset_ID := strconv.Itoa(Asset_ID)
	assetbytes,err := ctx.GetStub().GetState(asset_ID)
	if err != nil {
		return nil,err 
	}
	if assetbytes == nil {
		return nil, fmt.Errorf(`asset(%d) is not exists!!!`,Asset_ID)
	}
	var asset Assetstruct
	err = json.Unmarshal(assetbytes,&asset)
	if err != nil{
		return nil,err 
	}

    return &asset , nil
}


func (s *managingJongchin) DeleteAsset(ctx contractapi.TransactionContextInterface,Asset_ID int ) error {
	//DelState(key string) error
	exists, err := s.AssetExists(ctx,asset_ID)
	if err != nil {
		return err
	}
	if exists ==false {
		return fmt.Errorf(`AssetID (%d) does not exists`,asset_ID)
	} 

	asset_ID :=strconv.Itoa(Asset_ID)
	return ctx.GetStub().DelState(asset_ID)
}
func (s *managingJongchin) AddOwner(ctx contractapi.TransactionContextInterface, Asset_ID int, Owner string, share int) error{ 
	 
	asset, err := s.ReadAsset(ctx,Asset_ID)
	if err!=nil {
		return err
	}
	for _, owner := range asset.Owner {
		if owner.OwnerName == Owner{  
			return fmt.Errorf(`Owner already exists.`)
		}
	}
	var assetOwner = AssetOwner{OwnerName : Owner, Stake:share}
	asset.Owner = append(asset.Owner, assetOwner)

	assetbytes, err := json.Marshal(asset)
	if err != nil {
		return err
	}
	assetID := strconv.Itoa(Asset_ID)
	return ctx.GetStub().PutState(assetID, assetbytes)
}

// func (s *managingJongchin) DeleteOwner(ctx contractapi.TransactionContextInterface, Asset_ID int, Owner string) error{
// 	for _, owner := range asset.Owner {
// 		if owner != Owner{  
	
// }

func main() {

	chaincode, err := contractapi.NewChaincode(new(managingJongchin))

	if err != nil {
		fmt.Printf("Error create teamate chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting teamate chaincode: %s", err.Error())
	}
}
