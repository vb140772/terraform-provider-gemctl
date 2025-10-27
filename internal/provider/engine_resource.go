package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	client "github.com/vb140772/terraform-provider-gemctl/internal/client"
)

// Ensure NewEngineResource returns a resource with the correct interface implementation
var _ resource.Resource = &engineResource{}

type engineResource struct {
	client *client.GeminiClient
}

type engineResourceModel struct {
	ID           types.String `tfsdk:"id"`
	EngineID     types.String `tfsdk:"engine_id"`
	DisplayName  types.String `tfsdk:"display_name"`
	DataStores   types.List   `tfsdk:"data_stores"`
	Name         types.String `tfsdk:"name"`
}

func NewEngineResource(c *client.GeminiClient) resource.Resource {
	return &engineResource{
		client: c,
	}
}

func (r *engineResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_engine"
}

func (r *engineResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a search engine in Google Gemini Enterprise. An engine can be connected to multiple data stores to provide search capabilities.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"engine_id": schema.StringAttribute{
				Required:    true,
				Description: "Unique identifier for the engine",
			},
			"display_name": schema.StringAttribute{
				Required:    true,
				Description: "Display name for the engine",
			},
			"data_stores": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "List of data store IDs to connect to this engine",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "Full resource name of the engine",
			},
		},
	}
}

func (r *engineResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var model engineResourceModel
	diags := req.Plan.Get(ctx, &model)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert data stores list
	var dataStoreIDs []string
	if !model.DataStores.IsNull() {
		for _, ds := range model.DataStores.Elements() {
			dataStoreIDs = append(dataStoreIDs, ds.(types.String).ValueString())
		}
	}

	// Create the engine
	result, err := r.client.CreateSearchEngine(
		model.EngineID.ValueString(),
		model.DisplayName.ValueString(),
		dataStoreIDs,
		"SEARCH_TIER_ENTERPRISE",
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating engine",
			fmt.Sprintf("Failed to create engine: %v", err),
		)
		return
	}

	model.ID = types.StringValue(model.EngineID.ValueString())
	model.Name = types.StringValue(result.EngineName)

	diags = resp.State.Set(ctx, model)
	resp.Diagnostics.Append(diags...)
}

func (r *engineResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var model engineResourceModel
	diags := req.State.Get(ctx, &model)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the full engine name
	engineName := fmt.Sprintf("projects/%s/locations/%s/collections/%s/engines/%s",
		r.client.Config().ProjectID,
		r.client.Config().Location,
		r.client.Config().Collection,
		model.EngineID.ValueString())

	// Read the engine
	engine, err := r.client.GetEngineDetails(engineName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading engine",
			fmt.Sprintf("Failed to read engine: %v", err),
		)
		return
	}

	model.DisplayName = types.StringValue(engine.DisplayName)
	model.Name = types.StringValue(engine.Name)

	diags = resp.State.Set(ctx, model)
	resp.Diagnostics.Append(diags...)
}

func (r *engineResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var model engineResourceModel
	diags := req.Plan.Get(ctx, &model)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert data stores list
	var dataStoreIDs []string
	if !model.DataStores.IsNull() {
		for _, ds := range model.DataStores.Elements() {
			dataStoreIDs = append(dataStoreIDs, ds.(types.String).ValueString())
		}
	}

	// Update the engine (create new version with updated config)
	result, err := r.client.CreateSearchEngine(
		model.EngineID.ValueString(),
		model.DisplayName.ValueString(),
		dataStoreIDs,
		"SEARCH_TIER_ENTERPRISE",
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating engine",
			fmt.Sprintf("Failed to update engine: %v", err),
		)
		return
	}

	model.Name = types.StringValue(result.EngineName)

	diags = resp.State.Set(ctx, model)
	resp.Diagnostics.Append(diags...)
}

func (r *engineResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var model engineResourceModel
	diags := req.State.Get(ctx, &model)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the full engine name
	engineName := fmt.Sprintf("projects/%s/locations/%s/collections/%s/engines/%s",
		r.client.Config().ProjectID,
		r.client.Config().Location,
		r.client.Config().Collection,
		model.EngineID.ValueString())

	// Delete the engine
	_, err := r.client.DeleteEngine(engineName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting engine",
			fmt.Sprintf("Failed to delete engine: %v", err),
		)
		return
	}

	resp.State.RemoveResource(ctx)
}

