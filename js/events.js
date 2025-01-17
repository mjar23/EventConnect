// Event class representing an event
class Event {
    constructor(event) {
      this.id = event.id;
      this.eventname = event.eventname;
      this.venue = event.venue;
      this.date = event.date;
      this.imageurl = event.imageurl;
      this.weather = event.weather;
    }
  
    // Render the event card HTML
    renderCard() {
      const eventCard = document.createElement('div');
      eventCard.classList.add('event-card');
  
      const eventCardImage = document.createElement('div');
      eventCardImage.classList.add('event-card-image');
      const img = document.createElement('img');
      img.src = this.imageurl;
      img.alt = this.eventname;
      eventCardImage.appendChild(img);
  
      const eventCardContent = document.createElement('div');
      eventCardContent.classList.add('event-card-content');
      eventCardContent.innerHTML = `
        <h3 class="event-title"><i class="fas fa-calendar-alt"></i> ${this.eventname}</h3>
        <p class="event-venue"><i class="fas fa-map-marker-alt"></i> ${this.venue.name}</p>
        <p class="event-location"><i class="fas fa-map-pin"></i> ${this.venue.postcode_lookup}, ${this.venue.country}</p>
        <p class="event-date"><i class="fas fa-calendar-day"></i> ${this.date}</p>
        ${this.weather ? `<p class="event-weather"><i class="fas fa-cloud-rain"></i> Weather: ${this.weather}</p>` : ''}
      `;
  
      const eventCardFooter = document.createElement('div');
      eventCardFooter.classList.add('event-card-footer');
      const viewDetailsButton = document.createElement('button');
      viewDetailsButton.classList.add('button', 'is-primary');
      viewDetailsButton.textContent = 'View Details';
      viewDetailsButton.addEventListener('click', () => this.viewDetails());
      eventCardFooter.appendChild(viewDetailsButton);
  
      eventCard.appendChild(eventCardImage);
      eventCard.appendChild(eventCardContent);
      eventCard.appendChild(eventCardFooter);
  
      return eventCard;
    }
  
    // View event details
    viewDetails() {
      window.location.href = `/event-details.html?eventId=${this.id}`;
    }
  }
  
  // EventManager class responsible for managing events
  class EventManager {
    constructor() {
      this.events = [];
      this.map = null;
      this.userLatitude = null;
      this.userLongitude = null;
      
    }
  
    // Display events in the UI
    displayEvents() {
      const eventsContainer = document.getElementById('events-container');
      eventsContainer.innerHTML = '';
  
      this.events.forEach(event => {
        const eventCard = event.renderCard();
        eventsContainer.appendChild(eventCard);
      });
    }
  
    // Get weather information for an event
    async getWeather(latitude, longitude) {
      const apiKey = "65bd6689ba60f41871774e40059c6129";
      const url = `http://api.openweathermap.org/data/2.5/weather?lat=${latitude}&lon=${longitude}&appid=${apiKey}`;
  
      try {
        const response = await fetch(url);
        const weatherData = await response.json();
        console.log('Weather Data:', weatherData);
        const weather = `${weatherData.weather[0].main}, Temperature: ${(weatherData.main.temp - 273.15).toFixed(1)}°C`;
        console.log('Weather:', weather);
        return weather;
      } catch (error) {
        console.error('Error fetching weather:', error);
        return 'Weather information unavailable';
      }
    }
  
    // Fetch events from the server
    async fetchEvents(filterParams) {
      const queryString = new URLSearchParams(filterParams).toString();
      const url = `http://localhost:8000/events?${queryString}`;
  
      try {
        const response = await fetch(url);
        const eventsData = await response.json();
  
        // Create an array of promises for fetching weather data for each event
        const weatherPromises = eventsData.map(async (eventData) => {
          const weather = await this.getWeather(eventData.venue.latitude, eventData.venue.longitude);
          eventData.weather = weather;
          return new Event(eventData);
        });
  
        // Wait for all weather data to be fetched before updating events
        this.events = await Promise.all(weatherPromises);
  
        this.displayEvents();
        this.addEventsToMap();
      } catch (error) {
        console.error('Error fetching events:', error);
      }
    }
  
    // Initialize the Google Map
    initMap() {
      const mapOptions = {
        zoom: 8,
        center: { lat: this.userLatitude || 51.5072, lng: this.userLongitude || -0.1275 },
      };
  
      this.map = new google.maps.Map(document.getElementById('map'), mapOptions);
    }
  
    // event markers to the map
    addEventsToMap() {
      if (!this.map) return;
  
      const infoWindow = new google.maps.InfoWindow();
  
      this.events.forEach((event) => {
        const marker = new google.maps.Marker({
          position: { lat: event.venue.latitude, lng: event.venue.longitude },
          map: this.map,
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
          infoWindow.open(this.map, marker);
        });
  
        // Add mouseout event listener to the marker
        marker.addListener('mouseout', () => {
          infoWindow.close();
        });
      });
  
      // Center the map on the first event's location
      if (this.events.length > 0) {
        const latitude = this.events[0].venue.latitude;
        const longitude = this.events[0].venue.longitude;
        const center = { lat: latitude, lng: longitude };
        this.map.setCenter(center);
        this.map.setZoom(11);
      }
    }
  
    // Get user's location
    getLocation() {
      if (navigator.geolocation) {
        navigator.geolocation.getCurrentPosition(
          position => this.showPosition(position),
          error => this.showError(error)
        );
      } else {
        console.log("Geolocation is not supported by this browser.");
        this.initMap();
        this.fetchEvents();
      }
    }
  
    // Show user's position on the map
    showPosition(position) {
      this.userLatitude = position.coords.latitude;
      this.userLongitude = position.coords.longitude;
  
      this.initMap();
      this.fetchEvents();
    }
  
    // Handle geolocation errors
    showError(error) {
      console.error('Error getting location:', error);
      this.initMap();
      this.fetchEvents();
    }
  }
  
  // Initialize the EventManager
  const eventManager = new EventManager();
  
  // Add event listener to the filter form
  const filterForm = document.getElementById('filterForm');
  filterForm.addEventListener('submit', async (event) => {
    event.preventDefault();
    const formData = new FormData(event.target);
    const filterParams = {};
  
    // Handle location filter
    const location = formData.get('location');
    if (!location) {
      alert('Location is required.');
      return;
    }
    filterParams['country'] = 'GB';
    filterParams['keyword'] = location;
  
    // Handle event type filters
    const eventcodes = Array.from(formData.getAll('eventcode'), code => code.toUpperCase());
    if (eventcodes.length > 0) {
      filterParams['eventcode'] = eventcodes.join(',');
    }
  
    // Handle genre filters
    const genres = Array.from(formData.getAll('genre'), genre => genre.toUpperCase());
    if (genres.length > 0) {
      filterParams['g'] = genres.join(',');
    }
  
    // Handle other filters
    const keyword = formData.get('keyword');
    if (keyword) {
      filterParams['keyword'] = keyword;
    }
  
    const minDate = formData.get('minDate');
    const maxDate = formData.get('maxDate');
    if (minDate && maxDate && minDate > maxDate) {
      alert('Start date cannot be greater than end date.');
      return;
    }
    if (minDate) {
      filterParams['minDate'] = minDate;
    }
    if (maxDate) {
      filterParams['maxDate'] = maxDate;
    }
  
    try {
      await eventManager.fetchEvents(filterParams);
    } catch (error) {
      console.error('Error fetching events:', error);
    }
  });
  
  // event listeners to toggle filter and map visibility
  document.getElementById('toggle-filters').addEventListener('click', function() {
    const filterContainer = document.getElementById('filter-container');
    filterContainer.style.display = filterContainer.style.display === 'none' ? 'block' : 'none';
  });
  
  document.getElementById('toggle-map').addEventListener('click', function() {
    const mapContainer = document.getElementById('map-container');
    mapContainer.style.display = mapContainer.style.display === 'none' ? 'block' : 'none';
  });
  
  // Get user's location and initialize the map
  eventManager.getLocation();