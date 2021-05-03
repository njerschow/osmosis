package keeper

import (
	"encoding/json"
	"fmt"

	"github.com/c-osmosis/osmosis/x/incentives/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// getLastPotID returns ID used last time
func (k Keeper) getLastPotID(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.KeyLastPotID)
	if bz == nil {
		return 0
	}

	return sdk.BigEndianToUint64(bz)
}

// setLastPotID save ID used by last pot
func (k Keeper) setLastPotID(ctx sdk.Context, ID uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.KeyLastPotID, sdk.Uint64ToBigEndian(ID))
}

// potStoreKey returns action store key from ID
func potStoreKey(ID uint64) []byte {
	return combineKeys(types.KeyPrefixPeriodPot, sdk.Uint64ToBigEndian(ID))
}

// getPotRefs get pot IDs specified on the prefix and timestamp key
func (k Keeper) getPotRefs(ctx sdk.Context, key []byte) []uint64 {
	store := ctx.KVStore(k.storeKey)
	potIDs := []uint64{}
	if store.Has(key) {
		bz := store.Get(key)
		err := json.Unmarshal(bz, &potIDs)
		if err != nil {
			panic(err)
		}
	}
	return potIDs
}

// addPotRefByKey append pot ID into an array associated to provided key
func (k Keeper) addPotRefByKey(ctx sdk.Context, key []byte, potID uint64) error {
	store := ctx.KVStore(k.storeKey)
	potIDs := k.getPotRefs(ctx, key)
	if findIndex(potIDs, potID) > -1 {
		return fmt.Errorf("pot with same ID exist: %d", potID)
	}
	potIDs = append(potIDs, potID)
	bz, err := json.Marshal(potIDs)
	if err != nil {
		return err
	}
	store.Set(key, bz)
	return nil
}

// deletePotRefByKey removes pot ID from an array associated to provided key
func (k Keeper) deletePotRefByKey(ctx sdk.Context, key []byte, potID uint64) {
	var index = -1
	store := ctx.KVStore(k.storeKey)
	potIDs := k.getPotRefs(ctx, key)
	potIDs, index = removeValue(potIDs, potID)
	if index < 0 {
		panic(fmt.Sprintf("specific pot with ID %d not found", potID))
	}
	if len(potIDs) == 0 {
		store.Delete(key)
	} else {
		bz, err := json.Marshal(potIDs)
		if err != nil {
			panic(err)
		}
		store.Set(key, bz)
	}
}