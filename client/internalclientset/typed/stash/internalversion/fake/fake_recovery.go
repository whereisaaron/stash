/*
Copyright 2018 The Stash Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package fake

import (
	stash "github.com/appscode/stash/apis/stash"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeRecoveries implements RecoveryInterface
type FakeRecoveries struct {
	Fake *FakeStash
	ns   string
}

var recoveriesResource = schema.GroupVersionResource{Group: "stash.appscode.com", Version: "", Resource: "recoveries"}

var recoveriesKind = schema.GroupVersionKind{Group: "stash.appscode.com", Version: "", Kind: "Recovery"}

// Get takes name of the recovery, and returns the corresponding recovery object, and an error if there is any.
func (c *FakeRecoveries) Get(name string, options v1.GetOptions) (result *stash.Recovery, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(recoveriesResource, c.ns, name), &stash.Recovery{})

	if obj == nil {
		return nil, err
	}
	return obj.(*stash.Recovery), err
}

// List takes label and field selectors, and returns the list of Recoveries that match those selectors.
func (c *FakeRecoveries) List(opts v1.ListOptions) (result *stash.RecoveryList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(recoveriesResource, recoveriesKind, c.ns, opts), &stash.RecoveryList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &stash.RecoveryList{}
	for _, item := range obj.(*stash.RecoveryList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested recoveries.
func (c *FakeRecoveries) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(recoveriesResource, c.ns, opts))

}

// Create takes the representation of a recovery and creates it.  Returns the server's representation of the recovery, and an error, if there is any.
func (c *FakeRecoveries) Create(recovery *stash.Recovery) (result *stash.Recovery, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(recoveriesResource, c.ns, recovery), &stash.Recovery{})

	if obj == nil {
		return nil, err
	}
	return obj.(*stash.Recovery), err
}

// Update takes the representation of a recovery and updates it. Returns the server's representation of the recovery, and an error, if there is any.
func (c *FakeRecoveries) Update(recovery *stash.Recovery) (result *stash.Recovery, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(recoveriesResource, c.ns, recovery), &stash.Recovery{})

	if obj == nil {
		return nil, err
	}
	return obj.(*stash.Recovery), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeRecoveries) UpdateStatus(recovery *stash.Recovery) (*stash.Recovery, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(recoveriesResource, "status", c.ns, recovery), &stash.Recovery{})

	if obj == nil {
		return nil, err
	}
	return obj.(*stash.Recovery), err
}

// Delete takes name of the recovery and deletes it. Returns an error if one occurs.
func (c *FakeRecoveries) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(recoveriesResource, c.ns, name), &stash.Recovery{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeRecoveries) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(recoveriesResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &stash.RecoveryList{})
	return err
}

// Patch applies the patch and returns the patched recovery.
func (c *FakeRecoveries) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *stash.Recovery, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(recoveriesResource, c.ns, name, data, subresources...), &stash.Recovery{})

	if obj == nil {
		return nil, err
	}
	return obj.(*stash.Recovery), err
}
