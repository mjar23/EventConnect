// google/google-maps.js

let map;

function initMap(userLocations, userLatitude, userLongitude) {
    if (userLocations && userLocations.length > 0) {
        // Event details page with user locations
        const bounds = new google.maps.LatLngBounds();
        userLocations.forEach(location => {
            const latLng = new google.maps.LatLng(location.latitude, location.longitude);
            bounds.extend(latLng);
        });

        map = new google.maps.Map(document.getElementById('map'), {
            center: bounds.getCenter(),
            zoom: 8,
        });

        map.fitBounds(bounds);

        plotUserLocations(userLocations);
    } else {
        // Events page with user's current location or default location
        const mapOptions = {
            zoom: 8,
            center: { lat: userLatitude || 51.5072, lng: userLongitude || -0.1275 },
        };

        map = new google.maps.Map(document.getElementById('map'), mapOptions);
    }
}

function plotUserLocations(userLocations) {
    userLocations.forEach(location => {
        const latLng = new google.maps.LatLng(location.latitude, location.longitude);
        const marker = new google.maps.Marker({
            position: latLng,
            map: map,
            userId: location.id // Store the user ID as a data attribute
        });

        console.log('User ID:', location.id);

        // Add click event listener to the marker
        marker.addListener('click', () => {
            // Redirect to other-user-profile.html with the user ID as a query parameter
            const userId = marker.get('userId');
            window.location.href = `/other-user-profile.html?userId=${userId}`;
        });
    });
}

function addEventsToMap(events) {
    if (!map) return;

    const infoWindow = new google.maps.InfoWindow();

    events.forEach((event) => {
        const marker = new google.maps.Marker({
            position: { lat: event.venue.latitude, lng: event.venue.longitude },
            map: map,
            title: `${event.eventname} (${event.venue.name})`,
        });

        // Add click event listener to the marker
        marker.addListener('click', () => {
            event.viewDetails();
        });

        // Add mouseover event listener to the marker
        marker.addListener('mouseover', () => {
            const contentString = `
                <div>
                    <h4>${event.eventname}</h4>
                    <p>${event.venue.name}</p>
                    <p>${event.date}</p>
                    ${event.weather ? `<p>Weather: ${event.weather}</p>` : ''}
                </div>
            `;
            infoWindow.setContent(contentString);
            infoWindow.open(map, marker);
        });

        // Add mouseout event listener to the marker
        marker.addListener('mouseout', () => {
            infoWindow.close();
        });
    });

    // Center the map on the first event's location
    if (events.length > 0) {
        const latitude = events[0].venue.latitude;
        const longitude = events[0].venue.longitude;
        const center = { lat: latitude, lng: longitude };
        map.setCenter(center);
        map.setZoom(11);
    }
}

function loadGoogleMapsAPI(callback) {
    const script = document.createElement('script');
    script.src = `https://maps.googleapis.com/maps/api/js?key=AIzaSyBDR_N5soK76D55aICV-bW2jDEOxE5WwVQ&callback=${callback}`;
    script.async = true;
    script.defer = true;
    document.head.appendChild(script);
}

export { initMap, addEventsToMap, loadGoogleMapsAPI };