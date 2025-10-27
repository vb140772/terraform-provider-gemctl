package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	client "github.com/vb140772/terraform-provider-gemctl/internal/client"
)

// Ensure the implementation satisfies the expected interfaces
var _ provider.Provider = &gemctlProvider{}

type gemctlProvider struct {
	client *client.GeminiClient
}

type gemctlProviderModel struct {
	ProjectID          types.String `tfsdk:"project_id"`
	Location           types.String `tfsdk:"location"`
	Collection         types.String `tfsdk:"collection"`
	UseServiceAccount  types.Bool   `tfsdk:"use_service_account"`
}

func New() provider.Provider {
	return &gemctlProvider{}
}

func (p *gemctlProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "gemctl"
	resp.Version = "0.1.0"
}

func (p *gemctlProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `The gemctl provider manages Google Gemini Enterprise resources, including search engines and data stores. The provider requires Google Cloud credentials configured via Application Default Credentials (ADC) or user credentials.`,
		Attributes: map[string]schema.Attribute{
			"project_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Google Cloud project ID where resources will be created.",
			},
			"location": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Location for resources (e.g., `us`, `global`). Defaults to `us`.",
			},
			"collection": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Collection ID for organizing resources. Defaults to `default_collection`.",
			},
			"use_service_account": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Use service account credentials instead of user credentials. When false, uses `gcloud auth print-access-token`.",
			},
		},
	}
}

func (p *gemctlProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config gemctlProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	location := config.Location.ValueString()
	if location == "" {
		location = "us"
	}
	collection := config.Collection.ValueString()
	if collection == "" {
		collection = "default_collection"
	}

	clientConfig := &client.Config{
		ProjectID:         config.ProjectID.ValueString(),
		Location:          location,
		Collection:        collection,
		UseServiceAccount: config.UseServiceAccount.ValueBool(),
	}

	geminiClient, err := client.NewGeminiClient(clientConfig)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create Gemini client",
			err.Error(),
		)
		return
	}

	resp.DataSourceData = geminiClient
	resp.ResourceData = geminiClient
	
	// Store client in provider
	p.client = geminiClient
}

func (p *gemctlProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		func() resource.Resource { return NewEngineResource(p.client) },
		func() resource.Resource { return NewDataStoreResource(p.client) },
	}
}

func (p *gemctlProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		func() datasource.DataSource { return NewEngineDataSource(p.client) },
		func() datasource.DataSource { return NewDataStoreDataSource(p.client) },
	}
}

