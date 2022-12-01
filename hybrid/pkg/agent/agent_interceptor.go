package hybridagent

import (
	"context"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/hewlettpackard/hybrid/pkg/common"
	nodeattestorv1 "github.com/spiffe/spire-plugin-sdk/proto/spire/plugin/agent/nodeattestor/v1"
)

type AgentInterceptor interface {
	Recv() (*nodeattestorv1.Challenge, error)
	Send(challenge *nodeattestorv1.PayloadOrChallengeResponse) error
	Context() context.Context
	SetLogger(logger hclog.Logger)
	SetPluginName(name string)
	GetMessage() common.PluginMessage
	setCustomStream(stream nodeattestorv1.NodeAttestor_AidAttestationServer)
}

type HybridPluginAgentInterceptor struct {
	stream nodeattestorv1.NodeAttestor_AidAttestationServer
	nodeattestorv1.NodeAttestor_AidAttestationServer
	ctx        context.Context
	logger     hclog.Logger
	payload    []byte
	pluginName string
}

func (m *HybridPluginAgentInterceptor) SetPluginName(name string) {
	m.pluginName = name
}

func (m *HybridPluginAgentInterceptor) GetMessage() common.PluginMessage {
	return common.PluginMessage{
		PluginName: m.pluginName,
		PluginData: m.payload,
	}
}

func (m *HybridPluginAgentInterceptor) Recv() (*nodeattestorv1.Challenge, error) {
	return m.stream.Recv()
}

func (m *HybridPluginAgentInterceptor) Send(challenge *nodeattestorv1.PayloadOrChallengeResponse) error {
	m.payload = challenge.GetPayload()
	return nil
}

func (m *HybridPluginAgentInterceptor) Context() context.Context {
	return m.ctx
}

func (m *HybridPluginAgentInterceptor) SetLogger(logger hclog.Logger) {
	m.logger = logger
}

func (m *HybridPluginAgentInterceptor) setCustomStream(stream nodeattestorv1.NodeAttestor_AidAttestationServer) {
	m.stream = stream
	m.ctx = stream.Context()
}

func NewAgentInterceptor() AgentInterceptor {
	return &HybridPluginAgentInterceptor{}
}
