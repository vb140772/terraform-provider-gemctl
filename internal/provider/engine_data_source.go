package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	client "github.com/vb140772/terraform-provider-gemctl/internal/client"
)

type engineDataSource struct {
	client *client.GeminiClient
}

type engineDataSourceModel struct {
	EngineID    types.String `tfsdk:"engine_id"`
	Name        types.String `tfsdk:"name"`
	DisplayName types.String `tfsdk:"display_name"`
	SolutionType types.String `tfsdk:"solution_type"`
	IndustryVertical types.String `tfsdk:"industry_vertical"`
	DataStoreIds types.List `tfsdk:"data_store_ids"`
}

func NewEngineDataSource(c *client.GeminiClient) datasource.DataSource {
	return &engineDataSource{
		client: c,
	}
}

func (d *engineDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_engine"
}

func (d *engineDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"engine_id": schema.StringAttribute{
				Required:    true,
				Description: "Engine ID to look up",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "Full resource name of the engine",
			},
			"display_name": schema.StringAttribute{
				Computed:    true,
				Description: "Display name of the engine",
			},
			"solution_type": schema.StringAttribute{
				Computed:    true,
				Description: "Solution type of the engine",
			},
			"industry_vertical": schema.StringAttribute{
				Computed:    true,
				Description: "Industry vertical of the engine",
			},
			"data_store_ids": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "List of data store IDs connected to this engine",
			},
		},
	}
}

func (d *engineDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var model engineDataSourceModel
	diags := req.Config.Get(ctx, &model)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the full engine name
	engineName := fmt.Sprintf("projects/%s/locations/%s/collections/%s/engines/%s",
		d.client.Config().ProjectID,
		d.client.Config().Location,
		d.client.Config().Collection,
		model.EngineID.ValueString())

	// Read the engine
	engine, err := d.client.GetEngineDetails(engineName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading engine",
			fmt.Sprintf("Failed to read engine: %v", err),
		)
		return
	}

	model.Name = types.StringValue(engine.Name)
	model.DisplayName = types.StringValue(engine.DisplayName)
	model.SolutionType = types.StringValue(engine.SolutionType)
	model.IndustryVertical = types.StringValue(engine.IndustryVertical)

	// Convert data store IDs to list
	dataStoreList := []types.String{}
	for _, dsID := range engine.DataStoreIds {
		dataStoreList = append(dataStoreList, types.StringValue(dsID))
	}
	model.DataStoreIds, diags = types.ListValueFrom(ctx, types.StringType, dataStoreList)
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, model)
	resp.Diagnostics.Append(diags...)
}

