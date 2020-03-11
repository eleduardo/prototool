package lint

import (
	"errors"
	"plugin"

	"github.com/emicklei/proto"
	"github.com/uber/prototool/internal/settings"
	"github.com/uber/prototool/internal/text"
)

type validatorDelegate interface {
	GetId() string
	GetPurpose() string
	OnStart(proto *proto.Proto, workingdir string, dirpath string, FileData string) error
	Finally() error
	VisitMessage(message *proto.Message)
	VisitService(v *proto.Service)
	VisitSyntax(s *proto.Syntax)
	VisitPackage(p *proto.Package)
	VisitOption(o *proto.Option)
	VisitImport(i *proto.Import)
	VisitNormalField(i *proto.NormalField)
	VisitEnumField(i *proto.EnumField)
	VisitEnum(e *proto.Enum)
	VisitComment(e *proto.Comment)
	VisitOneof(o *proto.Oneof)
	VisitOneofField(o *proto.OneOfField)
	VisitReserved(r *proto.Reserved)
	VisitRPC(r *proto.RPC)
	VisitMapField(f *proto.MapField)
	VisitGroup(g *proto.Group)
	VisitExtensions(e *proto.Extensions)
}

type pluggableLinter struct {
	failures []*text.Failure
	delegate validatorDelegate
}

func (pl *pluggableLinter) VisitMessage(m *proto.Message) {
	pl.delegate.VisitMessage(m)
}

func (pl *pluggableLinter) VisitService(v *proto.Service) {
	pl.delegate.VisitService(v)
}

func (pl *pluggableLinter) VisitSyntax(s *proto.Syntax) {
	pl.delegate.VisitSyntax(s)
}

func (pl *pluggableLinter) VisitPackage(p *proto.Package) {
	pl.delegate.VisitPackage(p)
}

func (pl *pluggableLinter) VisitOption(o *proto.Option) {
	pl.delegate.VisitOption(o)
}

func (pl *pluggableLinter) VisitImport(i *proto.Import) {
	pl.delegate.VisitImport(i)
}

func (pl *pluggableLinter) VisitNormalField(i *proto.NormalField) {
	pl.delegate.VisitNormalField(i)
}

func (pl *pluggableLinter) VisitEnumField(i *proto.EnumField) {
	pl.delegate.VisitEnumField(i)
}

func (pl *pluggableLinter) VisitEnum(e *proto.Enum) {
	pl.delegate.VisitEnum(e)
}

func (pl *pluggableLinter) VisitComment(e *proto.Comment) {
	pl.delegate.VisitComment(e)
}

func (pl *pluggableLinter) VisitOneof(o *proto.Oneof) {
	pl.delegate.VisitOneof(o)
}

func (pl *pluggableLinter) VisitOneofField(o *proto.OneOfField) {
	pl.delegate.VisitOneofField(o)
}

func (pl *pluggableLinter) VisitReserved(r *proto.Reserved) {
	pl.delegate.VisitReserved(r)
}

func (pl *pluggableLinter) VisitRPC(r *proto.RPC) {
	pl.delegate.VisitRPC(r)
}

func (pl *pluggableLinter) VisitMapField(f *proto.MapField) {
	pl.delegate.VisitMapField(f)
}

func (pl *pluggableLinter) VisitGroup(g *proto.Group) {
	pl.delegate.VisitGroup(g)
}

func (pl *pluggableLinter) VisitExtensions(e *proto.Extensions) {
	pl.delegate.VisitExtensions(e)
}

func (pl *pluggableLinter) OnStart(d *FileDescriptor) error {
	return pl.delegate.OnStart(d.Proto, d.ProtoSet.WorkDirPath, d.ProtoSet.DirPath, d.FileData)
}

func (pl *pluggableLinter) Finally() error {
	return pl.delegate.Finally()
}

func (pl *pluggableLinter) ID() string {
	return pl.delegate.GetId()
}

func (pl *pluggableLinter) Purpose(config settings.LintConfig) string {
	return pl.delegate.GetPurpose()
}

func (pl *pluggableLinter) Check(dirPath string, descriptors []*FileDescriptor) ([]*text.Failure, error) {
	err := runVisitor(pl, descriptors)
	for _, failure := range pl.failures {
		failure.LintID = pl.delegate.GetId()
	}
	return pl.failures, err
}

func NewPluggableVisitor(modulelocation string, description string) (pluggableLinter, error) {
	plugin, err := plugin.Open(modulelocation)
	if err != nil {
		return pluggableLinter{}, err
	}
	//load the visitor delegate symbol
	delegatesymbol, err := plugin.Lookup("ValidatorDelegate")
	if err != nil {
		return pluggableLinter{}, err
	}

	var delegate validatorDelegate
	delegate, ok := delegatesymbol.(validatorDelegate)
	if !ok {
		return pluggableLinter{}, errors.New("Unable to cast plugin symbol")
	}
	return pluggableLinter{delegate: delegate}, nil
}
