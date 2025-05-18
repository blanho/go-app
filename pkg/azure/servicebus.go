// pkg/azure/servicebus.go
package azure

import (
	"context"
	"errors"
	"fmt"

	servicebus "github.com/Azure/azure-service-bus-go"
)

type ServiceBus struct {
	namespace *servicebus.Namespace
	queues    map[string]*servicebus.Queue
	topics    map[string]*servicebus.Topic
}

func NewServiceBus(connectionString string) (*ServiceBus, error) {
	namespace, err := servicebus.NewNamespace(servicebus.NamespaceWithConnectionString(connectionString))
	if err != nil {
		return nil, fmt.Errorf("failed to create Service Bus namespace: %w", err)
	}

	return &ServiceBus{
		namespace: namespace,
		queues:    make(map[string]*servicebus.Queue),
		topics:    make(map[string]*servicebus.Topic),
	}, nil
}

func (sb *ServiceBus) GetQueue(ctx context.Context, name string) (*servicebus.Queue, error) {
	if queue, exists := sb.queues[name]; exists {
		return queue, nil
	}

	queue, err := sb.namespace.NewQueue(name)
	if err != nil {
		return nil, fmt.Errorf("failed to create queue %s: %w", name, err)
	}

	sb.queues[name] = queue
	return queue, nil
}

func (sb *ServiceBus) SendMessage(ctx context.Context, queueName string, data []byte) error {
	queue, err := sb.GetQueue(ctx, queueName)
	if err != nil {
		return err
	}

	message := servicebus.Message{
		Data: data,
	}

	return queue.Send(ctx, &message)
}

func (sb *ServiceBus) Close(ctx context.Context) error {
	var errs []error

	for name, queue := range sb.queues {
		if err := queue.Close(ctx); err != nil {
			errs = append(errs, fmt.Errorf("failed to close queue %s: %w", name, err))
		}
	}

	for name, topic := range sb.topics {
		if err := topic.Close(ctx); err != nil {
			errs = append(errs, fmt.Errorf("failed to close topic %s: %w", name, err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors closing service bus: %v", errs)
	}

	return nil
}