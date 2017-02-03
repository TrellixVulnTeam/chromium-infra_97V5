// Code generated by protoc-gen-go.
// source: infra/tricium/api/v1/config.proto
// DO NOT EDIT!

/*
Package tricium is a generated protocol buffer package.

It is generated from these files:
	infra/tricium/api/v1/config.proto
	infra/tricium/api/v1/data.proto
	infra/tricium/api/v1/tricium.proto

It has these top-level messages:
	ServiceConfig
	ProjectDetails
	ProjectConfig
	RepoDetails
	GitRepoDetails
	Acl
	Platform
	Selection
	Analyzer
	ConfigDef
	Impl
	Recipe
	Config
	Cmd
	CipdPackage
	Data
	AnalyzeRequest
	AnalyzeResponse
*/
package tricium

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// Supported kinds of repositories.
type RepoDetails_Kind int32

const (
	RepoDetails_GIT RepoDetails_Kind = 0
)

var RepoDetails_Kind_name = map[int32]string{
	0: "GIT",
}
var RepoDetails_Kind_value = map[string]int32{
	"GIT": 0,
}

func (x RepoDetails_Kind) String() string {
	return proto.EnumName(RepoDetails_Kind_name, int32(x))
}
func (RepoDetails_Kind) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{3, 0} }

// Roles relevant to Tricium.
type Acl_Role int32

const (
	// Can read progress/results.
	Acl_READER Acl_Role = 0
	// Can request analysis.
	Acl_REQUESTER Acl_Role = 1
)

var Acl_Role_name = map[int32]string{
	0: "READER",
	1: "REQUESTER",
}
var Acl_Role_value = map[string]int32{
	"READER":    0,
	"REQUESTER": 1,
}

func (x Acl_Role) String() string {
	return proto.EnumName(Acl_Role_name, int32(x))
}
func (Acl_Role) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{5, 0} }

// Tricium service configuration.
//
// Listing supported platforms and analyzers shared between projects connected
// to Tricium.
type ServiceConfig struct {
	// Supported platforms.
	Platforms []*Platform `protobuf:"bytes,1,rep,name=platforms" json:"platforms,omitempty"`
	// List of shared analyzers.
	Analyzers []*Analyzer `protobuf:"bytes,2,rep,name=analyzers" json:"analyzers,omitempty"`
	// Details for connected projects.
	Projects []*ProjectDetails `protobuf:"bytes,3,rep,name=projects" json:"projects,omitempty"`
	// PubSub topic used to collect worker completion notifications from Swarming.
	SwarmingWorkerTopic string `protobuf:"bytes,4,opt,name=swarming_worker_topic,json=swarmingWorkerTopic" json:"swarming_worker_topic,omitempty"`
}

func (m *ServiceConfig) Reset()                    { *m = ServiceConfig{} }
func (m *ServiceConfig) String() string            { return proto.CompactTextString(m) }
func (*ServiceConfig) ProtoMessage()               {}
func (*ServiceConfig) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *ServiceConfig) GetPlatforms() []*Platform {
	if m != nil {
		return m.Platforms
	}
	return nil
}

func (m *ServiceConfig) GetAnalyzers() []*Analyzer {
	if m != nil {
		return m.Analyzers
	}
	return nil
}

func (m *ServiceConfig) GetProjects() []*ProjectDetails {
	if m != nil {
		return m.Projects
	}
	return nil
}

func (m *ServiceConfig) GetSwarmingWorkerTopic() string {
	if m != nil {
		return m.SwarmingWorkerTopic
	}
	return ""
}

type ProjectDetails struct {
	// Project name used to map these project details to the config for a project.
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	// General service account for this project.
	// Used for any service interaction, with the exception of swarming.
	ServiceAccount string `protobuf:"bytes,2,opt,name=service_account,json=serviceAccount" json:"service_account,omitempty"`
	// Project-specific swarming service account.
	SwarmingServiceAccount string `protobuf:"bytes,3,opt,name=swarming_service_account,json=swarmingServiceAccount" json:"swarming_service_account,omitempty"`
	// Details of the repository connected to the project. This should be the
	// repository hosting the files that should be analyzed for this project.
	RepoDetails *RepoDetails `protobuf:"bytes,4,opt,name=repo_details,json=repoDetails" json:"repo_details,omitempty"`
}

func (m *ProjectDetails) Reset()                    { *m = ProjectDetails{} }
func (m *ProjectDetails) String() string            { return proto.CompactTextString(m) }
func (*ProjectDetails) ProtoMessage()               {}
func (*ProjectDetails) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *ProjectDetails) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *ProjectDetails) GetServiceAccount() string {
	if m != nil {
		return m.ServiceAccount
	}
	return ""
}

func (m *ProjectDetails) GetSwarmingServiceAccount() string {
	if m != nil {
		return m.SwarmingServiceAccount
	}
	return ""
}

func (m *ProjectDetails) GetRepoDetails() *RepoDetails {
	if m != nil {
		return m.RepoDetails
	}
	return nil
}

// Tricium project configuration.
//
// Specifies details needed to connect a project to Tricium, adds project
// specific analyzers and implementations, and selects analyzer
// implementations.
type ProjectConfig struct {
	// Project name,
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	// Access control rules for the project.
	Acls []*Acl `protobuf:"bytes,2,rep,name=acls" json:"acls,omitempty"`
	// Project-specific analyzer details. This includes project-specific analyzer
	// implementations and full project-specific analyzer specifications.
	Analyzers []*Analyzer `protobuf:"bytes,3,rep,name=analyzers" json:"analyzers,omitempty"`
	// Selection of analyzer implementations to run for this project.
	Selections []*Selection `protobuf:"bytes,4,rep,name=selections" json:"selections,omitempty"`
}

func (m *ProjectConfig) Reset()                    { *m = ProjectConfig{} }
func (m *ProjectConfig) String() string            { return proto.CompactTextString(m) }
func (*ProjectConfig) ProtoMessage()               {}
func (*ProjectConfig) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *ProjectConfig) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *ProjectConfig) GetAcls() []*Acl {
	if m != nil {
		return m.Acls
	}
	return nil
}

func (m *ProjectConfig) GetAnalyzers() []*Analyzer {
	if m != nil {
		return m.Analyzers
	}
	return nil
}

func (m *ProjectConfig) GetSelections() []*Selection {
	if m != nil {
		return m.Selections
	}
	return nil
}

// Repository details for a project.
type RepoDetails struct {
	Kind RepoDetails_Kind `protobuf:"varint,1,opt,name=kind,enum=tricium.RepoDetails_Kind" json:"kind,omitempty"`
	// If repository kind is GIT then provide Git details.
	GitDetails *GitRepoDetails `protobuf:"bytes,2,opt,name=git_details,json=gitDetails" json:"git_details,omitempty"`
}

func (m *RepoDetails) Reset()                    { *m = RepoDetails{} }
func (m *RepoDetails) String() string            { return proto.CompactTextString(m) }
func (*RepoDetails) ProtoMessage()               {}
func (*RepoDetails) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *RepoDetails) GetKind() RepoDetails_Kind {
	if m != nil {
		return m.Kind
	}
	return RepoDetails_GIT
}

func (m *RepoDetails) GetGitDetails() *GitRepoDetails {
	if m != nil {
		return m.GitDetails
	}
	return nil
}

// Git repository details.
type GitRepoDetails struct {
	// URL to repository.
	Repository string `protobuf:"bytes,1,opt,name=repository" json:"repository,omitempty"`
	// Default ref to use to get files to analyze.
	Ref string `protobuf:"bytes,2,opt,name=ref" json:"ref,omitempty"`
}

func (m *GitRepoDetails) Reset()                    { *m = GitRepoDetails{} }
func (m *GitRepoDetails) String() string            { return proto.CompactTextString(m) }
func (*GitRepoDetails) ProtoMessage()               {}
func (*GitRepoDetails) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *GitRepoDetails) GetRepository() string {
	if m != nil {
		return m.Repository
	}
	return ""
}

func (m *GitRepoDetails) GetRef() string {
	if m != nil {
		return m.Ref
	}
	return ""
}

// Access control rules.
type Acl struct {
	// Role of a group or identity.
	Role Acl_Role `protobuf:"varint,1,opt,name=role,enum=tricium.Acl_Role" json:"role,omitempty"`
	// Name of group, as defined in the auth service. Specify either group or
	// identity, not both.
	Group string `protobuf:"bytes,2,opt,name=group" json:"group,omitempty"`
	// Identity, as defined by the auth service. Can be either an email address
	// or an indentity string, for instance, "anonymous:anonymous" for anonymous
	// users. Specify either group or identity, not both.
	Identity string `protobuf:"bytes,3,opt,name=identity" json:"identity,omitempty"`
}

func (m *Acl) Reset()                    { *m = Acl{} }
func (m *Acl) String() string            { return proto.CompactTextString(m) }
func (*Acl) ProtoMessage()               {}
func (*Acl) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *Acl) GetRole() Acl_Role {
	if m != nil {
		return m.Role
	}
	return Acl_READER
}

func (m *Acl) GetGroup() string {
	if m != nil {
		return m.Group
	}
	return ""
}

func (m *Acl) GetIdentity() string {
	if m != nil {
		return m.Identity
	}
	return ""
}

// Specification of a platform configuration.
type Platform struct {
	// Name use to refer to this platform configuration.
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	// Swarming dimensions of the form “key:value”, with keys and values mapping
	// to valid swarming keys/values.
	Dimensions []string `protobuf:"bytes,2,rep,name=dimensions" json:"dimensions,omitempty"`
}

func (m *Platform) Reset()                    { *m = Platform{} }
func (m *Platform) String() string            { return proto.CompactTextString(m) }
func (*Platform) ProtoMessage()               {}
func (*Platform) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *Platform) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Platform) GetDimensions() []string {
	if m != nil {
		return m.Dimensions
	}
	return nil
}

// Selection of analyzer implementations to run for a project.
type Selection struct {
	// Name of analyzer to run.
	Analyzer string `protobuf:"bytes,1,opt,name=analyzer" json:"analyzer,omitempty"`
	// Name of platform configuration to run analyzer on.
	Platform string `protobuf:"bytes,2,opt,name=platform" json:"platform,omitempty"`
	// Analyzer configuration to use on this platform.
	Configs []*Config `protobuf:"bytes,3,rep,name=configs" json:"configs,omitempty"`
}

func (m *Selection) Reset()                    { *m = Selection{} }
func (m *Selection) String() string            { return proto.CompactTextString(m) }
func (*Selection) ProtoMessage()               {}
func (*Selection) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *Selection) GetAnalyzer() string {
	if m != nil {
		return m.Analyzer
	}
	return ""
}

func (m *Selection) GetPlatform() string {
	if m != nil {
		return m.Platform
	}
	return ""
}

func (m *Selection) GetConfigs() []*Config {
	if m != nil {
		return m.Configs
	}
	return nil
}

// Analyzer specification.
type Analyzer struct {
	// Name of analyzer. This name is used to select the analyzer and is used
	// when reporting results for the analyzer. This name should be unique among
	// Tricium analyzers.
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	// Tricium data needed by this analyzer.
	Needs Data_Type `protobuf:"varint,2,opt,name=needs,enum=tricium.Data_Type" json:"needs,omitempty"`
	// Tricium data provided by this analyzer.
	Provides Data_Type `protobuf:"varint,3,opt,name=provides,enum=tricium.Data_Type" json:"provides,omitempty"`
	// Paths to run this analyzer on, defined as a glob.
	PathFilters []string `protobuf:"bytes,4,rep,name=path_filters,json=pathFilters" json:"path_filters,omitempty"`
	// Email to the owner of this analyzer.
	Owner string `protobuf:"bytes,6,opt,name=owner" json:"owner,omitempty"`
	// Monorail bug component for bug filing.
	Component string `protobuf:"bytes,7,opt,name=component" json:"component,omitempty"`
	// Analyzer configuration. These configuration options enable projects to
	// customize how an analyzer implementation analyzes their files.  It's
	// common for analyzers to provide a list of possible checks which can be
	// configured via a command line flag or similar. This field provides a way
	// to expose such flags as configuration options.
	ConfigDefs []*ConfigDef `protobuf:"bytes,8,rep,name=config_defs,json=configDefs" json:"config_defs,omitempty"`
	// Analyzer implementations. An analyzer may run on many platforms and this
	// may require many different implementations of the analyzer. An
	// implementation may be shared between several platforms if possible.
	Impls []*Impl `protobuf:"bytes,9,rep,name=impls" json:"impls,omitempty"`
}

func (m *Analyzer) Reset()                    { *m = Analyzer{} }
func (m *Analyzer) String() string            { return proto.CompactTextString(m) }
func (*Analyzer) ProtoMessage()               {}
func (*Analyzer) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

func (m *Analyzer) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Analyzer) GetNeeds() Data_Type {
	if m != nil {
		return m.Needs
	}
	return Data_NONE
}

func (m *Analyzer) GetProvides() Data_Type {
	if m != nil {
		return m.Provides
	}
	return Data_NONE
}

func (m *Analyzer) GetPathFilters() []string {
	if m != nil {
		return m.PathFilters
	}
	return nil
}

func (m *Analyzer) GetOwner() string {
	if m != nil {
		return m.Owner
	}
	return ""
}

func (m *Analyzer) GetComponent() string {
	if m != nil {
		return m.Component
	}
	return ""
}

func (m *Analyzer) GetConfigDefs() []*ConfigDef {
	if m != nil {
		return m.ConfigDefs
	}
	return nil
}

func (m *Analyzer) GetImpls() []*Impl {
	if m != nil {
		return m.Impls
	}
	return nil
}

// Definition of an analyzer configuration, e.g., ClangTidy is configured with
// a ‘checks’ flag.
type ConfigDef struct {
	// Name of configuration option.
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	// Default value for the config, e.g., checks=”all”.
	Default string `protobuf:"bytes,2,opt,name=default" json:"default,omitempty"`
}

func (m *ConfigDef) Reset()                    { *m = ConfigDef{} }
func (m *ConfigDef) String() string            { return proto.CompactTextString(m) }
func (*ConfigDef) ProtoMessage()               {}
func (*ConfigDef) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

func (m *ConfigDef) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *ConfigDef) GetDefault() string {
	if m != nil {
		return m.Default
	}
	return ""
}

// Analyzer implementation for one or more platforms. Implementation can be
// recipe-based or binary-based.
type Impl struct {
	// Platforms this implementation applies to.
	Platforms []string `protobuf:"bytes,1,rep,name=platforms" json:"platforms,omitempty"`
	// Cipd packages needed by this implementation.
	CipdPackages []*CipdPackage `protobuf:"bytes,2,rep,name=cipd_packages,json=cipdPackages" json:"cipd_packages,omitempty"`
	// Recipe for recipe-based implementation. Either recipe or cmd, not both.
	Recipe *Recipe `protobuf:"bytes,3,opt,name=recipe" json:"recipe,omitempty"`
	// Command for binary-based implementation. Either recipe or cmd, not both.
	Cmd *Cmd `protobuf:"bytes,4,opt,name=cmd" json:"cmd,omitempty"`
	// Deadline for execution of corresponding worker (in minutes). Note that
	// this deadline includes the launch of a swarming task for the corresponding
	// worker, and collection of results from that worker.
	Deadline int32 `protobuf:"varint,5,opt,name=deadline" json:"deadline,omitempty"`
}

func (m *Impl) Reset()                    { *m = Impl{} }
func (m *Impl) String() string            { return proto.CompactTextString(m) }
func (*Impl) ProtoMessage()               {}
func (*Impl) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{10} }

func (m *Impl) GetPlatforms() []string {
	if m != nil {
		return m.Platforms
	}
	return nil
}

func (m *Impl) GetCipdPackages() []*CipdPackage {
	if m != nil {
		return m.CipdPackages
	}
	return nil
}

func (m *Impl) GetRecipe() *Recipe {
	if m != nil {
		return m.Recipe
	}
	return nil
}

func (m *Impl) GetCmd() *Cmd {
	if m != nil {
		return m.Cmd
	}
	return nil
}

func (m *Impl) GetDeadline() int32 {
	if m != nil {
		return m.Deadline
	}
	return 0
}

// Specification of how to find a recipe.
type Recipe struct {
	// Repository URL of the recipe package.
	Repository string `protobuf:"bytes,1,opt,name=repository" json:"repository,omitempty"`
	// Path to recipe in the repository.
	Path string `protobuf:"bytes,2,opt,name=path" json:"path,omitempty"`
	// Revision to use.
	Revision string `protobuf:"bytes,3,opt,name=revision" json:"revision,omitempty"`
	// A json string containing properties for the recipe.
	Properties string `protobuf:"bytes,4,opt,name=properties" json:"properties,omitempty"`
}

func (m *Recipe) Reset()                    { *m = Recipe{} }
func (m *Recipe) String() string            { return proto.CompactTextString(m) }
func (*Recipe) ProtoMessage()               {}
func (*Recipe) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{11} }

func (m *Recipe) GetRepository() string {
	if m != nil {
		return m.Repository
	}
	return ""
}

func (m *Recipe) GetPath() string {
	if m != nil {
		return m.Path
	}
	return ""
}

func (m *Recipe) GetRevision() string {
	if m != nil {
		return m.Revision
	}
	return ""
}

func (m *Recipe) GetProperties() string {
	if m != nil {
		return m.Properties
	}
	return ""
}

// Analyzer configuration used when selecting an analyzer implementation.
type Config struct {
	// Name of the configuration option.
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	// Value of the configuration.
	Value string `protobuf:"bytes,2,opt,name=value" json:"value,omitempty"`
}

func (m *Config) Reset()                    { *m = Config{} }
func (m *Config) String() string            { return proto.CompactTextString(m) }
func (*Config) ProtoMessage()               {}
func (*Config) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{12} }

func (m *Config) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Config) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

// Specification of a command.
type Cmd struct {
	// Executable binary.
	Exec string `protobuf:"bytes,1,opt,name=exec" json:"exec,omitempty"`
	// Arguments in order.
	Args []string `protobuf:"bytes,2,rep,name=args" json:"args,omitempty"`
}

func (m *Cmd) Reset()                    { *m = Cmd{} }
func (m *Cmd) String() string            { return proto.CompactTextString(m) }
func (*Cmd) ProtoMessage()               {}
func (*Cmd) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{13} }

func (m *Cmd) GetExec() string {
	if m != nil {
		return m.Exec
	}
	return ""
}

func (m *Cmd) GetArgs() []string {
	if m != nil {
		return m.Args
	}
	return nil
}

// CIPD package.
type CipdPackage struct {
	// CIPD package name.
	PackageName string `protobuf:"bytes,1,opt,name=package_name,json=packageName" json:"package_name,omitempty"`
	// Path to directory, relative to the working directory, where to install
	// package. Cannot be empty or start with a slash.
	Path string `protobuf:"bytes,2,opt,name=path" json:"path,omitempty"`
	// Package version.
	Version string `protobuf:"bytes,3,opt,name=version" json:"version,omitempty"`
}

func (m *CipdPackage) Reset()                    { *m = CipdPackage{} }
func (m *CipdPackage) String() string            { return proto.CompactTextString(m) }
func (*CipdPackage) ProtoMessage()               {}
func (*CipdPackage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{14} }

func (m *CipdPackage) GetPackageName() string {
	if m != nil {
		return m.PackageName
	}
	return ""
}

func (m *CipdPackage) GetPath() string {
	if m != nil {
		return m.Path
	}
	return ""
}

func (m *CipdPackage) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}

func init() {
	proto.RegisterType((*ServiceConfig)(nil), "tricium.ServiceConfig")
	proto.RegisterType((*ProjectDetails)(nil), "tricium.ProjectDetails")
	proto.RegisterType((*ProjectConfig)(nil), "tricium.ProjectConfig")
	proto.RegisterType((*RepoDetails)(nil), "tricium.RepoDetails")
	proto.RegisterType((*GitRepoDetails)(nil), "tricium.GitRepoDetails")
	proto.RegisterType((*Acl)(nil), "tricium.Acl")
	proto.RegisterType((*Platform)(nil), "tricium.Platform")
	proto.RegisterType((*Selection)(nil), "tricium.Selection")
	proto.RegisterType((*Analyzer)(nil), "tricium.Analyzer")
	proto.RegisterType((*ConfigDef)(nil), "tricium.ConfigDef")
	proto.RegisterType((*Impl)(nil), "tricium.Impl")
	proto.RegisterType((*Recipe)(nil), "tricium.Recipe")
	proto.RegisterType((*Config)(nil), "tricium.Config")
	proto.RegisterType((*Cmd)(nil), "tricium.Cmd")
	proto.RegisterType((*CipdPackage)(nil), "tricium.CipdPackage")
	proto.RegisterEnum("tricium.RepoDetails_Kind", RepoDetails_Kind_name, RepoDetails_Kind_value)
	proto.RegisterEnum("tricium.Acl_Role", Acl_Role_name, Acl_Role_value)
}

func init() { proto.RegisterFile("infra/tricium/api/v1/config.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 921 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x8c, 0x55, 0x51, 0x6f, 0xe3, 0x44,
	0x10, 0x3e, 0x37, 0x4e, 0xd2, 0x4c, 0x9a, 0x5c, 0x6f, 0x29, 0x60, 0x4e, 0xa8, 0xb4, 0x46, 0xe8,
	0xca, 0xc3, 0x25, 0x22, 0x7d, 0xe0, 0xee, 0x05, 0xa9, 0xb4, 0xe1, 0x74, 0x42, 0x42, 0xc7, 0xb6,
	0x88, 0x37, 0xa2, 0x65, 0x3d, 0x09, 0x4b, 0x6d, 0xef, 0x6a, 0xed, 0xa4, 0x57, 0xde, 0x78, 0xe2,
	0xb7, 0xc0, 0xaf, 0xe0, 0x81, 0x5f, 0xc1, 0xaf, 0x41, 0xbb, 0xde, 0xb5, 0x9d, 0x2a, 0x3a, 0xf1,
	0xb6, 0x33, 0xdf, 0xcc, 0x78, 0xe7, 0xfb, 0x66, 0xc7, 0x70, 0x2a, 0xf2, 0xa5, 0x66, 0xd3, 0x52,
	0x0b, 0x2e, 0xd6, 0xd9, 0x94, 0x29, 0x31, 0xdd, 0x7c, 0x31, 0xe5, 0x32, 0x5f, 0x8a, 0xd5, 0x44,
	0x69, 0x59, 0x4a, 0xd2, 0x77, 0xe0, 0xd3, 0x4f, 0x76, 0xc6, 0x26, 0xac, 0x64, 0x55, 0x64, 0xfc,
	0x6f, 0x00, 0xa3, 0x6b, 0xd4, 0x1b, 0xc1, 0xf1, 0xd2, 0x56, 0x20, 0x53, 0x18, 0xa8, 0x94, 0x95,
	0x4b, 0xa9, 0xb3, 0x22, 0x0a, 0x4e, 0x3a, 0x67, 0xc3, 0xd9, 0x93, 0x89, 0x2b, 0x30, 0x79, 0xe3,
	0x10, 0xda, 0xc4, 0x98, 0x04, 0x96, 0xb3, 0xf4, 0xfe, 0x37, 0xd4, 0x45, 0xb4, 0xf7, 0x20, 0xe1,
	0xc2, 0x21, 0xb4, 0x89, 0x21, 0xe7, 0xb0, 0xaf, 0xb4, 0xfc, 0x15, 0x79, 0x59, 0x44, 0x1d, 0x1b,
	0xff, 0x61, 0xf3, 0x81, 0x0a, 0xb8, 0xc2, 0x92, 0x89, 0xb4, 0xa0, 0x75, 0x20, 0x99, 0xc1, 0xfb,
	0xc5, 0x1d, 0xd3, 0x99, 0xc8, 0x57, 0x8b, 0x3b, 0xa9, 0x6f, 0x51, 0x2f, 0x4a, 0xa9, 0x04, 0x8f,
	0xc2, 0x93, 0xe0, 0x6c, 0x40, 0xdf, 0xf3, 0xe0, 0x8f, 0x16, 0xbb, 0x31, 0x50, 0xfc, 0x77, 0x00,
	0xe3, 0xed, 0x82, 0x84, 0x40, 0x98, 0xb3, 0x0c, 0xa3, 0xc0, 0x66, 0xd9, 0x33, 0x79, 0x06, 0x8f,
	0x8b, 0x8a, 0x82, 0x05, 0xe3, 0x5c, 0xae, 0xf3, 0x32, 0xda, 0xb3, 0xf0, 0xd8, 0xb9, 0x2f, 0x2a,
	0x2f, 0x79, 0x01, 0x51, 0x7d, 0x87, 0x87, 0x19, 0x1d, 0x9b, 0xf1, 0x81, 0xc7, 0xaf, 0xb7, 0x33,
	0xbf, 0x84, 0x03, 0x8d, 0x4a, 0x2e, 0x92, 0xea, 0x1a, 0xf6, 0xd2, 0xc3, 0xd9, 0x51, 0xdd, 0x36,
	0x45, 0x25, 0x7d, 0xcf, 0x43, 0xdd, 0x18, 0xf1, 0x5f, 0x01, 0x8c, 0x5c, 0x0b, 0x4e, 0x9f, 0x5d,
	0x1d, 0x9c, 0x40, 0xc8, 0x78, 0xea, 0xd9, 0x3f, 0x68, 0xd8, 0xe7, 0x29, 0xb5, 0xc8, 0xb6, 0x48,
	0x9d, 0xff, 0x21, 0xd2, 0x0c, 0xa0, 0xc0, 0x14, 0x79, 0x29, 0x64, 0x6e, 0xee, 0x6b, 0x32, 0x48,
	0x9d, 0x71, 0xed, 0x21, 0xda, 0x8a, 0x8a, 0xff, 0x08, 0x60, 0xd8, 0xea, 0x84, 0x3c, 0x87, 0xf0,
	0x56, 0xe4, 0x89, 0xbd, 0xea, 0x78, 0xf6, 0xd1, 0xae, 0x6e, 0x27, 0xdf, 0x8a, 0x3c, 0xa1, 0x36,
	0x8c, 0xbc, 0x80, 0xe1, 0x4a, 0x94, 0x35, 0x47, 0x7b, 0x96, 0xa3, 0x66, 0x34, 0x5e, 0x89, 0xb2,
	0x4d, 0x13, 0xac, 0x84, 0x57, 0x35, 0x7e, 0x0c, 0xa1, 0xa9, 0x43, 0xfa, 0xd0, 0x79, 0xf5, 0xfa,
	0xe6, 0xf0, 0x51, 0xfc, 0x35, 0x8c, 0xb7, 0xc3, 0xc9, 0x31, 0x80, 0xe1, 0xb5, 0x10, 0xa5, 0xd4,
	0xf7, 0x8e, 0xbc, 0x96, 0x87, 0x1c, 0x42, 0x47, 0xe3, 0xd2, 0x09, 0x6f, 0x8e, 0xf1, 0xef, 0x01,
	0x74, 0x2e, 0x78, 0x4a, 0x3e, 0x83, 0x50, 0xcb, 0x14, 0x5d, 0x17, 0x4f, 0xda, 0xe4, 0x4e, 0xa8,
	0x4c, 0x91, 0x5a, 0x98, 0x1c, 0x41, 0x77, 0xa5, 0xe5, 0x5a, 0xb9, 0x12, 0x95, 0x41, 0x9e, 0xc2,
	0xbe, 0x48, 0x30, 0x2f, 0x45, 0x79, 0xef, 0x46, 0xa4, 0xb6, 0xe3, 0x53, 0x08, 0x4d, 0x3e, 0x01,
	0xe8, 0xd1, 0xf9, 0xc5, 0xd5, 0x9c, 0x1e, 0x3e, 0x22, 0x23, 0x18, 0xd0, 0xf9, 0xf7, 0x3f, 0xcc,
	0xaf, 0x6f, 0xe6, 0xf4, 0x30, 0x88, 0xbf, 0x82, 0x7d, 0xff, 0xe4, 0x76, 0x0a, 0x7f, 0x0c, 0x90,
	0x88, 0x0c, 0xf3, 0xc2, 0xaa, 0x64, 0xe4, 0x1f, 0xd0, 0x96, 0x27, 0xce, 0x61, 0x50, 0x4b, 0x65,
	0xee, 0xe2, 0xf5, 0x75, 0x45, 0x6a, 0xdb, 0x60, 0xfe, 0x45, 0xbb, 0x06, 0x6a, 0x9b, 0x7c, 0x0e,
	0xfd, 0x6a, 0xbb, 0xf8, 0xc9, 0x79, 0x5c, 0x73, 0x50, 0xcd, 0x24, 0xf5, 0x78, 0xfc, 0xe7, 0x1e,
	0xec, 0xfb, 0x69, 0xda, 0x79, 0xe1, 0x33, 0xe8, 0xe6, 0x88, 0x49, 0xa5, 0xee, 0xb8, 0x35, 0x51,
	0x57, 0x66, 0x27, 0xdd, 0xdc, 0x2b, 0xa4, 0x55, 0x00, 0x99, 0xd8, 0x2d, 0xb1, 0x11, 0x09, 0x16,
	0x96, 0xb9, 0xdd, 0xc1, 0x75, 0x0c, 0x39, 0x85, 0x03, 0xc5, 0xca, 0x5f, 0x16, 0x4b, 0x91, 0x96,
	0x66, 0xc8, 0x43, 0x4b, 0xc6, 0xd0, 0xf8, 0xbe, 0xa9, 0x5c, 0x46, 0x22, 0x79, 0x97, 0xa3, 0x8e,
	0x7a, 0x95, 0x44, 0xd6, 0x20, 0x1f, 0xc3, 0x80, 0xcb, 0x4c, 0xc9, 0x1c, 0xf3, 0x32, 0xea, 0x5b,
	0xa4, 0x71, 0x90, 0x73, 0x18, 0x56, 0xcd, 0x2d, 0x12, 0x5c, 0x16, 0xd1, 0xfe, 0x83, 0x87, 0x50,
	0x11, 0x70, 0x85, 0x4b, 0x0a, 0xdc, 0x1f, 0x0b, 0xf2, 0x29, 0x74, 0x45, 0xa6, 0xd2, 0x22, 0x1a,
	0xd8, 0xf0, 0x51, 0x1d, 0xfe, 0x3a, 0x53, 0x29, 0xad, 0xb0, 0xf8, 0x25, 0x0c, 0xea, 0xec, 0x9d,
	0x5c, 0x45, 0xd0, 0x4f, 0x70, 0xc9, 0xd6, 0xa9, 0xdf, 0x47, 0xde, 0x8c, 0xff, 0x09, 0x20, 0x34,
	0xa5, 0xcc, 0xdd, 0xb7, 0x97, 0xf5, 0xa0, 0xbd, 0x99, 0x5f, 0xc2, 0x88, 0x0b, 0x95, 0x2c, 0x14,
	0xe3, 0xb7, 0x6c, 0x85, 0x7e, 0x3f, 0x34, 0x6b, 0xe7, 0x52, 0xa8, 0xe4, 0x4d, 0x05, 0xd2, 0x03,
	0xde, 0x18, 0x05, 0x79, 0x06, 0x3d, 0x8d, 0x5c, 0x28, 0xb4, 0xdc, 0xb7, 0x25, 0xa7, 0xd6, 0x4d,
	0x1d, 0x4c, 0x8e, 0xa1, 0xc3, 0xb3, 0xc4, 0x2d, 0xb4, 0x66, 0xf3, 0x5c, 0x66, 0x09, 0x35, 0x80,
	0x19, 0xac, 0x04, 0x59, 0x92, 0x8a, 0x1c, 0xa3, 0xee, 0x49, 0x70, 0xd6, 0xa5, 0xb5, 0x1d, 0xbf,
	0x85, 0x1e, 0xf5, 0x55, 0xde, 0xfd, 0x3a, 0x09, 0x84, 0x46, 0x48, 0xc7, 0x83, 0x3d, 0x9b, 0xca,
	0x1a, 0x37, 0xc2, 0x0c, 0xba, 0x7f, 0x5a, 0xde, 0x36, 0xf5, 0x94, 0x96, 0x0a, 0x75, 0x29, 0xb0,
	0x70, 0xbf, 0x88, 0x96, 0x27, 0x9e, 0x41, 0xef, 0x1d, 0xeb, 0xf4, 0x08, 0xba, 0x1b, 0x96, 0xae,
	0xd1, 0x3f, 0x65, 0x6b, 0xc4, 0xcf, 0xa1, 0x73, 0x99, 0x25, 0x26, 0x01, 0xdf, 0x22, 0xf7, 0x09,
	0xe6, 0x6c, 0x7c, 0x4c, 0xaf, 0xfc, 0x03, 0xb4, 0xe7, 0xf8, 0x27, 0x18, 0xb6, 0xe8, 0xad, 0xc6,
	0xd3, 0x1e, 0x17, 0xad, 0xef, 0x0d, 0x9d, 0xef, 0x3b, 0xf3, 0xd9, 0x5d, 0x4d, 0x46, 0xd0, 0xdf,
	0xa0, 0x6e, 0xf5, 0xe8, 0xcd, 0x9f, 0x7b, 0xf6, 0x07, 0x7e, 0xfe, 0x5f, 0x00, 0x00, 0x00, 0xff,
	0xff, 0x02, 0xec, 0x1e, 0x9e, 0x0f, 0x08, 0x00, 0x00,
}
