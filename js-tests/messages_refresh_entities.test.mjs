import {test} from 'node:test';
import assert from 'node:assert';

import {ServerRefreshDriverOffers, ServerRefreshPassengerInvitations} from '../assets/js/entities/messages.js';
import {Offer} from '../assets/js/entities/offer.js';
import {Invitation} from '../assets/js/entities/invitation.js';

test('ServerRefreshDriverOffers.getOffers() extracts Offer entities', () => {
    const message = new ServerRefreshDriverOffers().fromObject({
        offers: [
            {invitation_id: 'i1', status: 'pending', passenger_id: 'p1', lat: '54.68', lon: '25.28'},
            {invitation_id: 'i2', status: 'accepted', passenger_id: 'p2', lat: '54.70', lon: '25.30'},
        ],
    });

    const offers = message.offers

    assert.strictEqual(offers.length, 2);
    for (const offer of offers) {
        assert.ok(offer instanceof Offer);
    }

    assert.strictEqual(offers[0].getPassengerId(), 'p1');
    assert.strictEqual(offers[0].getStatus(), 'pending');
    assert.strictEqual(offers[0].getLat(), '54.68');
    assert.strictEqual(offers[0].getLng(), '25.28');

    assert.strictEqual(offers[1].getPassengerId(), 'p2');
    assert.strictEqual(offers[1].getStatus(), 'accepted');
    assert.strictEqual(offers[1].getLat(), '54.70');
    assert.strictEqual(offers[1].getLng(), '25.30');
});

test('ServerRefreshPassengerInvitations.getInvitations() extracts Invitation entities', () => {
    const message = new ServerRefreshPassengerInvitations().fromObject({
        invitations: [
            {invitation_id: 'i1', status: 'pending', driver_id: 'd1', lat: '54.68', lon: '25.28'},
            {invitation_id: 'i2', status: 'rejected', driver_id: 'd2', lat: '54.70', lon: '25.30'},
        ],
    });

    const invitations = message.getInvitations();

    assert.strictEqual(invitations.length, 2);
    for (const invitation of invitations) {
        assert.ok(invitation instanceof Invitation);
    }

    assert.strictEqual(invitations[0].getDriverId(), 'd1');
    assert.strictEqual(invitations[0].getStatus(), 'pending');
    assert.strictEqual(invitations[0].getLat(), '54.68');
    assert.strictEqual(invitations[0].getLon(), '25.28');

    assert.strictEqual(invitations[1].getDriverId(), 'd2');
    assert.strictEqual(invitations[1].getStatus(), 'rejected');
    assert.strictEqual(invitations[1].getLat(), '54.70');
    assert.strictEqual(invitations[1].getLon(), '25.30');
});
