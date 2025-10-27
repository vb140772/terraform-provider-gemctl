package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	client "github.com/vb140772/terraform-provider-gemctl/internal/client"
)

type dataStoreResource struct {
	client *client.GeminiClient
}

type dataStoreResourceModel struct {
	ID           types.String `tfsdk:"id"`
	DataStoreID  types.String `tfsdk:"data_store_id"`
	DisplayName  types.String `tfsdk:"display_name"`
	GCSUri       types.String `tfsdk:"gcs_uri"`
	Name         types.String `tfsdk:"name"`
}

func NewDataStoreResource(c *client.GeminiClient) resource.Resource {
	return &dataStoreResource{
		client: c,
	}
}

func (r *dataStoreResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_data_store"
}

func (r *dataStoreResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a data store in Google Gemini Enterprise. Data stores import content from GCS buckets and can be connected to search engines.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"data_store_id": schema.StringAttribute{
				Required:    true,
				Description: "Unique identifier for the data store",
			},
			"display_name": schema.StringAttribute{
				Required:    true,
				Description: "Display name for the data store",
			},
			"gcs_uri": schema.StringAttribute{
				Required:    true,
				Description: "GCS URI to import data from (e.g., gs://bucket/path/*)",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "Full resource name of the data store",
			},
		},
	}
}

func (r *dataStoreResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var model dataStoreResourceModel
	diags := req.Plan.Get(ctx, &model)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create data store from GCS
	result, err := r.client.CreateDataStoreFromGCS(
		model.DataStoreID.ValueString(),
		model.DisplayName.ValueString(),
		model.GCSUri.ValueString(),
		"DATA_SCHEMA_DOCUMENT", // Default schema
		"INCREMENTAL",           // Default reconciliation mode
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating data store",
			fmt.Sprintf("Failed to create data store: %v", err),
		)
		return
	}

	model.ID = types.StringValue(model.DataStoreID.ValueString())
	model.Name = types.StringValue(result.DataStoreName)

	diags = resp.State.Set(ctx, model)
	resp.Diagnostics.Append(diags...)
}

func (r *dataStoreResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var model dataStoreResourceModel
	diags := req.State.Get(ctx, &model)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the full data store name
	dataStoreName := fmt.Sprintf("projects/%s/locations/%s/collections/%s/dataStores/%s",
		r.client.Config().ProjectID,
		r.client.Config().Location,
		r.client.Config().Collection,
		model.DataStoreID.ValueString())

	// Read the data store
	dataStore, err := r.client.GetDataStoreDetails(dataStoreName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading data store",
			fmt.Sprintf("Failed to read data store: %v", err),
		)
		return
	}

	model.DisplayName = types.StringValue(dataStore.DisplayName)
	model.Name = types.StringValue(dataStore.Name)

	diags = resp.State.Set(ctx, model)
	resp.Diagnostics.Append(diags...)
}

func (r *dataStoreResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var model dataStoreResourceModel
	diags := req.Plan.Get(ctx, &model)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// For now, update recreates the data store with new config
	// In a real implementation, you might want to check what changed
	result, err := r.client.CreateDataStoreFromGCS(
		model.DataStoreID.ValueString(),
		model.DisplayName.ValueString(),
		model.GCSUri.ValueString(),
		"DATA_SCHEMA_DOCUMENT", // Default schema
		"INCREMENTAL",           // Default reconciliation mode
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating data store",
			fmt.Sprintf("Failed to update data store: %v", err),
		)
		return
	}

	model.Name = types.StringValue(result.DataStoreName)

	diags = resp.State.Set(ctx, model)
	resp.Diagnostics.Append(diags...)
}

func (r *dataStoreResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var model dataStoreResourceModel
	diags := req.State.Get(ctx, &model)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the full data store name
	dataStoreName := fmt.Sprintf("projects/%s/locations/%s/collections/%s/dataStores/%s",
		r.client.Config().ProjectID,
		r.client.Config().Location,
		r.client.Config().Collection,
		model.DataStoreID.ValueString())

	// Delete the data store
	_, err := r.client.DeleteDataStore(dataStoreName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting data store",
			fmt.Sprintf("Failed to delete data store: %v", err),
		)
		return
	}

	resp.State.RemoveResource(ctx)
}

