package templates

const ArgsTemplate string = `namespace {{.namespace}};

using MediatR;

public class {{.type}} : IRequest<Unit> { }
`

const HandlerTemplate string = `namespace {{.namespace}};

using Blockchain.NFT.Core.Repositories;
using MediatR;

public class {{.handler}} : IRequestHandler<{{.type}}, Unit>
{
    readonly IUnitOfWork unitOfWork;

    public {{.handler}}(IUnitOfWork unitOfWork)
    {
        this.unitOfWork = unitOfWork;
    }

    public Task<Unit> Handle({{.type}} {{.typeName}}, CancellationToken cancellationToken)
    {
        throw new NotImplementedException();
    }
}
`
