// Mirrors internal/messages/message_types.go — keep the two in sync.
export const MessageType = {
    // === to driver or to passenger
    SERVER_REQUESTS_COORDINATES: 'server_requests_coordinates',
    SERVER_REQUESTS_CLIENT_INFO: 'server_requests_client_info',

    // === driver or passenger
    CLIENT_RESPONDS_COORDINATES: 'client_responds_coordinates',
    CLIENT_RESPONDS_INFO: 'client_responds_info',
    CLIENT_RESPONSE_ERROR: 'client_response_error',
    CLIENT_ERROR: 'client_error',

    // --- driver-related messages
    SERVER_OFFERS_PASSENGER: 'server_offers_passenger',
    SERVER_CANCELS_OFFER: 'server_cancels_offer',
    DRIVER_CANCELS_OFFER: 'driver_cancels_offer',
    DRIVER_ACCEPTS_OFFER: 'driver_accepts_offer',
    DRIVER_REJECTS_OFFER: 'driver_rejects_offer',
    DRIVER_CHANGES_STATUS: 'driver_changes_status',

    // --- passenger-related messages
    PASSENGER_INVITES_DRIVER: 'passenger_invites_driver',
    PASSENGER_CANCELS_INVITE: 'passenger_cancels_invite',
    SERVER_NOTIFIES_INVITE_ACCEPTED: 'server_notifies_invite_accepted',
    SERVER_NOTIFIES_INVITE_REJECTED: 'server_notifies_invite_rejected',
    SERVER_NOTIFIES_INVITE_CANCELED: 'server_notifies_invite_canceled',
    PASSENGER_REQUEST_DRIVER_COORDS: 'passenger_request_driver_coordinates',
    SERVER_RESPONDS_DRIVER_COORDS: 'server_responds_driver_coordinates',
    SERVER_NOTIFIES_VOYAGE_STARTED: 'server_notifies_voyage_started',
    SERVER_NOTIFIES_VOYAGE_FINISHED: 'server_notifies_voyage_finished',

    SERVER_REFRESH_DRIVER_STATUS: 'server_refresh_driver_status',
    SERVER_REFRESH_DRIVER_OFFERS: 'server_refresh_driver_offers',
    SERVER_REFRESH_PASSENGER_INVITATIONS: 'server_refresh_passenger_invitations',
};

export const DriverStatus = {
    IDLE: 'idle',
    WORKING: 'working',
    RESTING: 'resting',
    OFFLINE: 'offline',
}

// used for invites and offers. Mirrors internal/state/invitation.go — keep the two in sync.
export const InvitationStatus = {
    PENDING: 'pending',
    ACCEPTED: 'accepted',
    REJECTED: 'rejected',
    CANCELED: 'cancelled',
    DRIVING: 'driving',
    COMPLETED: 'completed',
}
