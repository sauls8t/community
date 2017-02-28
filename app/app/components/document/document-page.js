// Copyright 2016 Documize Inc. <legal@documize.com>. All rights reserved.
//
// This software (Documize Community Edition) is licensed under
// GNU AGPL v3 http://www.gnu.org/licenses/agpl-3.0.en.html
//
// You can operate outside the AGPL restrictions by purchasing
// Documize Enterprise Edition and obtaining a commercial license
// by contacting <sales@documize.com>.
//
// https://documize.com

import Ember from 'ember';
import NotifierMixin from '../../mixins/notifier';
import TooltipMixin from '../../mixins/tooltip';

export default Ember.Component.extend(NotifierMixin, TooltipMixin, {
	documentService: Ember.inject.service('document'),
	sectionService: Ember.inject.service('section'),
	editMode: false,

	actions: {
		onAddBlock(block) {
			this.attrs.onAddBlock(block);
		},

		onCopyPage(documentId) {
			this.attrs.onCopyPage(this.get('page.id'), documentId);
		},

		onMovePage(documentId) {
			this.attrs.onMovePage(this.get('page.id'), documentId);
		},

		onDeletePage(deleteChildren) {
			let page = this.get('page');

			if (is.undefined(page)) {
				return;
			}

			let params = {
				id: page.get('id'),
				title: page.get('title'),
				children: deleteChildren
			};

			this.attrs.onDeletePage(params);
		},

		onEdit() {
			if (this.get('editMode')) {
				return;
			}

			let page = this.get('page');

			this.get('documentService').getPageMeta(page.get('documentId'), page.get('id')).then((meta) => {
				this.set('editMode', true);
				this.set('meta', meta);
			});
		},

		onCancelEdit() {
			this.set('editMode', false);
		},

		onSavePage(page, meta) {
			this.set('editMode', false);
			this.attrs.onSavePage(page, meta);
		}
	}
});
