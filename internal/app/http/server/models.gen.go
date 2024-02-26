// Package server provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.0.0 DO NOT EDIT.
package server

// Defines values for GetCompletionsParamsKind.
const (
	FetchNames GetCompletionsParamsKind = "fetchNames"
	Names      GetCompletionsParamsKind = "names"
)

// AuthorsResponse defines model for AuthorsResponse.
type AuthorsResponse struct {
	Data []string `json:"data"`
}

// ErrorResponse defines model for ErrorResponse.
type ErrorResponse struct {
	Code     *string `json:"code,omitempty"`
	Detail   *string `json:"detail,omitempty"`
	Instance *string `json:"instance,omitempty"`
	Status   int     `json:"status"`
	Title    string  `json:"title"`
	Type     *string `json:"type,omitempty"`
}

// InventoryEntry defines model for InventoryEntry.
type InventoryEntry struct {
	Links              *InventoryEntryLinks    `json:"links,omitempty"`
	Name               string                  `json:"name"`
	SchemaAuthor       SchemaAuthor            `json:"schema:author"`
	SchemaManufacturer SchemaManufacturer      `json:"schema:manufacturer"`
	SchemaMpn          string                  `json:"schema:mpn"`
	Versions           []InventoryEntryVersion `json:"versions"`
}

// InventoryEntryLinks defines model for InventoryEntryLinks.
type InventoryEntryLinks struct {
	Self string `json:"self"`
}

// InventoryEntryResponse defines model for InventoryEntryResponse.
type InventoryEntryResponse struct {
	Data InventoryEntry `json:"data"`
}

// InventoryEntryVersion defines model for InventoryEntryVersion.
type InventoryEntryVersion struct {
	Description string                      `json:"description"`
	Digest      string                      `json:"digest"`
	ExternalID  string                      `json:"externalID"`
	Links       *InventoryEntryVersionLinks `json:"links,omitempty"`
	Timestamp   string                      `json:"timestamp"`
	TmID        string                      `json:"tmID"`
	Version     ModelVersion                `json:"version"`
}

// InventoryEntryVersionLinks defines model for InventoryEntryVersionLinks.
type InventoryEntryVersionLinks struct {
	Content string `json:"content"`
}

// InventoryEntryVersionsResponse defines model for InventoryEntryVersionsResponse.
type InventoryEntryVersionsResponse struct {
	Data []InventoryEntryVersion `json:"data"`
}

// InventoryResponse defines model for InventoryResponse.
type InventoryResponse struct {
	Data []InventoryEntry `json:"data"`
	Meta *Meta            `json:"meta,omitempty"`
}

// ManufacturersResponse defines model for ManufacturersResponse.
type ManufacturersResponse struct {
	Data []string `json:"data"`
}

// Meta defines model for Meta.
type Meta struct {
	Page *MetaPage `json:"page,omitempty"`
}

// MetaPage defines model for MetaPage.
type MetaPage struct {
	Elements int `json:"elements"`
}

// ModelVersion defines model for ModelVersion.
type ModelVersion struct {
	Model string `json:"model"`
}

// MpnsResponse defines model for MpnsResponse.
type MpnsResponse struct {
	Data []string `json:"data"`
}

// PushThingModelResponse defines model for PushThingModelResponse.
type PushThingModelResponse struct {
	Data PushThingModelResult `json:"data"`
}

// PushThingModelResult defines model for PushThingModelResult.
type PushThingModelResult struct {
	TmID string `json:"tmID"`
}

// SchemaAuthor defines model for SchemaAuthor.
type SchemaAuthor struct {
	SchemaName string `json:"schema:name"`
}

// SchemaManufacturer defines model for SchemaManufacturer.
type SchemaManufacturer struct {
	SchemaName string `json:"schema:name"`
}

// GetCompletionsParams defines parameters for GetCompletions.
type GetCompletionsParams struct {
	// Kind Kind of data to complete
	Kind *GetCompletionsParamsKind `form:"kind,omitempty" json:"kind,omitempty"`

	// ToComplete Data to complete
	ToComplete *string `form:"toComplete,omitempty" json:"toComplete,omitempty"`
}

// GetCompletionsParamsKind defines parameters for GetCompletions.
type GetCompletionsParamsKind string

// GetAuthorsParams defines parameters for GetAuthors.
type GetAuthorsParams struct {
	// FilterManufacturer Filters the authors according to whether they have inventory entries
	// which belong to at least one of the given manufacturers with an exact match.
	// The filter works additive to other filters.
	FilterManufacturer *string `form:"filter.manufacturer,omitempty" json:"filter.manufacturer,omitempty"`

	// FilterMpn Filters the authors according to whether they have inventory entries
	// which belong to at least one of the given mpn (manufacturer part number) with an exact match.
	// The filter works additive to other filters.
	FilterMpn *string `form:"filter.mpn,omitempty" json:"filter.mpn,omitempty"`

	// Search Filters the authors according to whether they have inventory entries
	// where their content matches the given search.
	// The search works additive to other filters.
	Search *string `form:"search,omitempty" json:"search,omitempty"`
}

// GetInventoryParams defines parameters for GetInventory.
type GetInventoryParams struct {
	// FilterAuthor Filters the inventory by one or more authors having exact match.
	// The filter works additive to other filters.
	FilterAuthor *string `form:"filter.author,omitempty" json:"filter.author,omitempty"`

	// FilterManufacturer Filters the inventory by one or more manufacturers having exact match.
	// The filter works additive to other filters.
	FilterManufacturer *string `form:"filter.manufacturer,omitempty" json:"filter.manufacturer,omitempty"`

	// FilterMpn Filters the inventory by one ore more mpn (manufacturer part number) having exact match.
	// The filter works additive to other filters.
	FilterMpn *string `form:"filter.mpn,omitempty" json:"filter.mpn,omitempty"`

	// FilterName Filters the inventory by inventory entry name having a prefix match of full path parts.
	// The filter works additive to other filters.
	FilterName *string `form:"filter.name,omitempty" json:"filter.name,omitempty"`

	// Search Filters the inventory according to whether the content of the inventory entries matches the given search.
	// The search works additive to other filters.
	Search *string `form:"search,omitempty" json:"search,omitempty"`
}

// GetManufacturersParams defines parameters for GetManufacturers.
type GetManufacturersParams struct {
	// FilterAuthor Filters the manufacturers according to whether they belong to at least one of the given authors with an exact match.
	// The filter works additive to other filters.
	FilterAuthor *string `form:"filter.author,omitempty" json:"filter.author,omitempty"`

	// FilterMpn Filters the manufacturers according to whether they have inventory entries
	// which belong to at least one of the given mpn (manufacturer part number) with an exact match.
	// The filter works additive to other filters.
	FilterMpn *string `form:"filter.mpn,omitempty" json:"filter.mpn,omitempty"`

	// Search Filters the manufacturers according to whether they have inventory entries
	// where their content matches the given search.
	// The search works additive to other filters.
	Search *string `form:"search,omitempty" json:"search,omitempty"`
}

// GetMpnsParams defines parameters for GetMpns.
type GetMpnsParams struct {
	// FilterAuthor Filters the mpns according to whether they belong to at least one of the given authors with an exact match.
	// The filter works additive to other filters.
	FilterAuthor *string `form:"filter.author,omitempty" json:"filter.author,omitempty"`

	// FilterManufacturer Filters the mpns according to whether they belong to at least one of the given manufacturers with an exact match.
	// The filter works additive to other filters.
	FilterManufacturer *string `form:"filter.manufacturer,omitempty" json:"filter.manufacturer,omitempty"`

	// Search Filters the mpns according to whether their inventory entry content matches the given search.
	// The search works additive to other filters.
	Search *string `form:"search,omitempty" json:"search,omitempty"`
}

// PushThingModelJSONBody defines parameters for PushThingModel.
type PushThingModelJSONBody = map[string]interface{}

// GetThingModelByIdParams defines parameters for GetThingModelById.
type GetThingModelByIdParams struct {
	// RestoreId restore the TM's original external id, if it had one
	RestoreId *bool `form:"restoreId,omitempty" json:"restoreId,omitempty"`
}

// PushThingModelJSONRequestBody defines body for PushThingModel for application/json ContentType.
type PushThingModelJSONRequestBody = PushThingModelJSONBody
