package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	client "github.com/vb140772/terraform-provider-gemctl/internal/client"
)

type dataStoreDataSource struct {
	client *client.GeminiClient
}

type dataStoreDataSourceModel struct {
	DataStoreID      types.String `tfsdk:"data_store_id"`
	Name             types.String `tfsdk:"name"`
	DisplayName      types.String `tfsdk:"display_name"`
	IndustryVertical types.String `tfsdk:"industry_vertical"`
	ContentConfig    types.String `tfsdk:"content_config"`
	CreateTime       types.String `tfsdk:"create_time"`
}

func NewDataStoreDataSource(c *client.GeminiClient) datasource.DataSource {
	return &dataStoreDataSource{
		client: c,
	}
}

func (d *dataStoreDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_data_store"
}

func (d *dataStoreDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"data_store_id": schema.StringAttribute{
				Required:    true,
				Description: "Data store ID to look up",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "Full resource name of the data store",
			},
			"display_name": schema.StringAttribute{
				Computed:    true,
				Description: "Display name of the data store",
			},
			"industry_vertical": schema.StringAttribute{
				Computed:    true,
				Description: "Industry vertical of the data store",
			},
			"content_config": schema.StringAttribute{
				Computed:    true,
				Description: "Content configuration of the data store",
			},
			"create_time": schema.StringAttribute{
				Computed:    true,
				Description: "Creation time of the data store",
			},
		},
	}
}

func (d *dataStoreDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var model dataStoreDataSourceModel
	diags := req.Config.Get(ctx, &model)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the full data store name
	dataStoreName := fmt.Sprintf("projects/%s/locations/%s/collections/%s/dataStores/%s",
		d.client.Config().ProjectID,
		d.client.Config().Location,
		d.client.Config().Collection,
		model.DataStoreID.ValueString())

	// Read the data store
	dataStore, err := d.client.GetDataStoreDetails(dataStoreName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading data store",
			fmt.Sprintf("Failed to read data store: %v", err),
		)
		return
	}

	model.Name = types.StringValue(dataStore.Name)
	model.DisplayName = types.StringValue(dataStore.DisplayName)
	model.IndustryVertical = types.StringValue(dataStore.IndustryVertical)
	model.ContentConfig = types.StringValue(dataStore.ContentConfig)
	model.CreateTime = types.StringValue(dataStore.CreateTime)

	diags = resp.State.Set(ctx, model)
	resp.Diagnostics.Append(diags...)
}

