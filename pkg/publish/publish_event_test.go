package publish_test

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"fmt"
	"testing"

	"github.com/bhojpur/cms/pkg/publish"
	orm "github.com/bhojpur/orm/pkg/engine"
)

type createResourcePublishInterface struct {
}

func (createResourcePublishInterface) Publish(db *orm.DB, event publish.PublishEventInterface) error {
	if event, ok := event.(*publish.PublishEvent); ok {
		var product Product
		db.Set("publish:draft_mode", true).First(&product, event.Argument)
		pb.Publish(&product)
	}
	return nil
}

func (createResourcePublishInterface) Discard(db *orm.DB, event publish.PublishEventInterface) error {
	if event, ok := event.(*publish.PublishEvent); ok {
		var product Product
		db.Set("publish:draft_mode", true).First(&product, event.Argument)
		pb.Discard(&product)
	}
	return nil
}

type publishAllResourcesInterface struct {
}

func (publishAllResourcesInterface) Publish(db *orm.DB, event publish.PublishEventInterface) error {
	return nil
}

func (publishAllResourcesInterface) Discard(db *orm.DB, event publish.PublishEventInterface) error {
	return nil
}

func init() {
	publish.RegisterEvent("create_product", createResourcePublishInterface{})
	publish.RegisterEvent("publish_all_resources", publishAllResourcesInterface{})
}

func TestCreateNewEvent(t *testing.T) {
	product1 := Product{Name: "event_1"}
	pbdraft.Set("publish:publish_event", true).Save(&product1)
	event := publish.PublishEvent{Name: "create_product", Argument: fmt.Sprintf("%v", product1.ID)}
	db.Save(&event)

	if !pbprod.First(&Product{}, "name = ?", product1.Name).RecordNotFound() {
		t.Errorf("created resource in draft db with event should not be published to production db")
	}

	var productDraft Product
	if pbdraft.First(&productDraft, "name = ?", product1.Name).RecordNotFound() {
		t.Errorf("created resource in draft db with event should exist in draft db")
	}

	if productDraft.PublishStatus == publish.DIRTY {
		t.Errorf("product's publish status should not be DIRTY before publish event")
	}

	var publishEvent publish.PublishEvent
	if pbdraft.First(&publishEvent, "name = ?", "create_product").Error != nil {
		t.Errorf("created resource in draft db with event should create the event in db")
	}

	if !pbprod.First(&Product{}, "name = ?", product1.Name).RecordNotFound() {
		t.Errorf("product should not be published to production db before publish event")
	}

	publishEvent.Publish(db)

	if pbprod.First(&Product{}, "name = ?", product1.Name).RecordNotFound() {
		t.Errorf("product should be published to production db after publish event")
	}
}

func TestCreateProductWithPublishAllEvent(t *testing.T) {
	product1 := Product{Name: "event_1"}
	event := &publish.PublishEvent{Name: "publish_all_resources", Argument: "products"}
	pbdraft.Set("publish:publish_event", event).Save(&product1)
}
