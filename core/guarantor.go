package core

import (
	"github.com/pkg/errors"
	"log"

	utils "github.com/Varunram/essentials/utils"
	consts "github.com/YaleOpenLab/opensolar/consts"
	xlm "github.com/YaleOpenLab/openx/chains/xlm"
	assets "github.com/YaleOpenLab/openx/chains/xlm/assets"
	wallet "github.com/YaleOpenLab/openx/chains/xlm/wallet"
)

// AddFirstLossGuarantee adds the given entity as a first loss guarantor
func (a *Entity) AddFirstLossGuarantee(seedpwd string, amount float64) error {
	if !a.Guarantor {
		log.Println("caller not guarantor")
		return errors.New("caller not guarantor, quitting")
	}

	a.FirstLossGuarantee = seedpwd
	a.FirstLossGuaranteeAmt = amount
	return a.Save()
}

func (a *Entity) RefillEscrowAsset(projIndex int, asset string, amount float64, seedpwd string) error {
	if !a.Guarantor {
		log.Println("caller not guarantor")
		return errors.New("caller not guarantor, quitting")
	}

	project, err := RetrieveProject(projIndex)
	if err != nil {
		log.Println(err)
		return err
	}

	balancex, err := xlm.GetAssetBalance(a.U.StellarWallet.PublicKey, asset)
	if err != nil {
		log.Println(err)
		return err
	}

	balance, err := utils.ToFloat(balancex)
	if err != nil {
		log.Println(err)
		return err
	}

	if balance < amount {
		log.Println("guarantor does not required amount, refilling what amount they have")
		amount = balance - 1.0 // fees
	}

	seed, err := wallet.DecryptSeed(a.U.StellarWallet.EncryptedSeed, seedpwd)
	if err != nil {
		log.Println(err)
		return err
	}

	_, txhash, err := assets.SendAsset(consts.StablecoinCode, consts.StablecoinPublicKey,
		project.EscrowPubkey, amount, seed, "guarantor refund")
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("txhash: ", txhash)
	return nil
}

func (a *Entity) RefillEscrowXLM(projIndex int, amount float64, seedpwd string) error {
	if !a.Guarantor {
		log.Println("caller not guarantor")
		return errors.New("caller not guarantor, quitting")
	}

	project, err := RetrieveProject(projIndex)
	if err != nil {
		log.Println(err)
		return err
	}

	balancex, err := xlm.GetNativeBalance(a.U.StellarWallet.PublicKey)
	if err != nil {
		log.Println(err)
		return err
	}

	balance, err := utils.ToFloat(balancex)
	if err != nil {
		log.Println(err)
		return err
	}

	if balance < amount {
		log.Println("guarantor does not required amount, refilling what amount they have")
		amount = balance - 1.0 // fees
	}

	seed, err := wallet.DecryptSeed(a.U.StellarWallet.EncryptedSeed, seedpwd)
	if err != nil {
		log.Println(err)
		return err
	}

	_, txhash, err := xlm.SendXLM(project.EscrowPubkey, amount, seed, "guarantor refund")
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("txhash: ", txhash)
	return nil
}
