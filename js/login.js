class LoginForm {
  constructor() {
    this.formElement = document.getElementById('user-form');
    this.toggleFormLink = document.getElementById('toggle-form-link');
    this.formTitleElement = document.getElementById('form-title');
    this.formSubmitElement = document.getElementById('form-submit');
    this.messageElement = document.getElementById('message');
    this.formFields = this.getAllFormFields();

    this.initializeEventListeners();
    this.initializeLocationAutocomplete();
  }

  initializeEventListeners() {
    this.toggleFormLink.addEventListener('click', this.toggleForm.bind(this));
    this.formElement.addEventListener('submit', this.handleFormSubmit.bind(this));
  }

  initializeLocationAutocomplete() {
    const locationInput = document.getElementById('location');
    const autocomplete = new google.maps.places.Autocomplete(locationInput, {
      types: ['(cities)'],
      componentRestrictions: { country: 'GB' },
    });

    autocomplete.addListener('place_changed', () => {
      const place = autocomplete.getPlace();
      if (place.geometry) {
        document.getElementById('latitude').value = place.geometry.location.lat();
        document.getElementById('longitude').value = place.geometry.location.lng();
      }
    });
  }

  toggleForm(event) {
    event.preventDefault();
    this.isLoginForm() ? this.showSignupForm() : this.showLoginForm();
  }

  showSignupForm() {
    this.formTitleElement.textContent = 'Sign Up';
    this.formSubmitElement.textContent = 'Sign Up';
    this.showFormFields();
  }

  showLoginForm() {
    this.formTitleElement.textContent = 'Login';
    this.formSubmitElement.textContent = 'Login';
    this.hideFormFields();
  }

  getAllFormFields() {
    return [
      { id: 'email-field', type: 'input' },
      { id: 'first-name-field', type: 'input' },
      { id: 'last-name-field', type: 'input' },
      { id: 'bio-field', type: 'textarea' },
      { id: 'interests-field', type: 'input' },
      { id: 'location-field', type: 'input' },
      { id: 'age-field', type: 'input' },
      { id: 'gender-field', type: 'input' },
      { id: 'instagram-username-field', type: 'input' },
      { id: 'facebook-username-field', type: 'input' },
      { id: 'snapchat-username-field', type: 'input' }
    ];
  }

  showFormFields() {
    this.formFields.forEach(field => this.showField(field));
  }

  hideFormFields() {
    this.formFields.forEach(field => this.hideField(field));
  }

  showField(field) {
    const fieldElement = document.getElementById(field.id);
    fieldElement.style.display = 'block';
  }

  hideField(field) {
    const fieldElement = document.getElementById(field.id);
    fieldElement.style.display = 'none';
  }

  handleFormSubmit(event) {
    event.preventDefault();
    if (this.validateForm()) {
      this.submitForm();
    }
  }

  validateForm() {
    this.resetErrorMessages();
    let isValid = this.validateUsername() && this.validatePassword();
  
    if (this.isSignupForm()) {
      isValid = isValid && this.validateEmail() && this.validateFirstName() && this.validateLastName() && this.validateBio() && this.validateInterests() && this.validateLocation() && this.validateAge() && this.validateGender();
    }
  
    return isValid;
  }

  resetErrorMessages() {
    this.setErrorMessage('username-error', '');
    this.setErrorMessage('email-error', '');
    this.setErrorMessage('password-error', '');
    this.setErrorMessage('firstName-error', '');
    this.setErrorMessage('lastName-error', '');
    this.setErrorMessage('bio-error', '');
    this.setErrorMessage('interests-error', '');
    this.setErrorMessage('location-error', '');
    this.setErrorMessage('age-error', '');
    this.setErrorMessage('gender-error', '');
  }

  setErrorMessage(id, message) {
    const element = document.getElementById(id);
    element.textContent = message;
  }

  validateUsername() {
    const username = document.getElementById('username').value;
    if (!username) {
      this.setErrorMessage('username-error', 'Username is required.');
      return false;
    }
    return true;
  }

  validatePassword() {
    const password = document.getElementById('password').value;
    if (!password) {
      this.setErrorMessage('password-error', 'Password is required.');
      return false;
    }
    return true;
  }

  validateEmail() {
    const email = document.getElementById('email').value;
    if (!email) {
      this.setErrorMessage('email-error', 'Email is required.');
      return false;
    } else if (!this.isValidEmail(email)) {
      this.setErrorMessage('email-error', 'Invalid email format.');
      return false;
    }
    return true;
  }

  isValidEmail(email) {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return emailRegex.test(email);
  }

  validateFirstName() {
    const firstName = document.getElementById('firstName').value;
    if (!firstName) {
      this.setErrorMessage('firstName-error', 'First name is required.');
      return false;
    }
    return true;
  }

  validateLastName() {
    const lastName = document.getElementById('lastName').value;
    if (!lastName) {
      this.setErrorMessage('lastName-error', 'Last name is required.');
      return false;
    }
    return true;
  }

  validateBio() {
    const bio = document.getElementById('bio').value;
    if (!bio) {
      this.setErrorMessage('bio-error', 'Bio is required.');
      return false;
    }
    return true;
  }

  validateInterests() {
    const interests = document.getElementById('interests').value;
    if (!interests) {
      this.setErrorMessage('interests-error', 'Interests are required.');
      return false;
    }
    return true;
  }

  validateLocation() {
    const location = document.getElementById('location').value;
    if (!location) {
      this.setErrorMessage('location-error', 'Location is required.');
      return false;
    }
    return true;
  }

  validateAge() {
    const age = document.getElementById('age').value;
    if (!age) {
      this.setErrorMessage('age-error', 'Age is required.');
      return false;
    }
    return true;
  }

  validateGender() {
    const gender = document.getElementById('gender').value;
    if (!gender) {
      this.setErrorMessage('gender-error', 'Gender is required.');
      return false;
    }
    return true;
  }

  isSignupForm() {
    return this.formTitleElement.textContent === 'Sign Up';
  }

  isLoginForm() {
    return this.formTitleElement.textContent === 'Login';
  }

  submitForm() {
    const requestBody = this.getRequestBody();
    const endpoint = this.getEndpoint();
    const successMessage = this.getSuccessMessage();

    this.sendRequest(requestBody, endpoint, successMessage);
  }

  getRequestBody() {
    const requestBody = {
      username: document.getElementById('username').value,
      password: document.getElementById('password').value
    };

    if (this.isSignupForm()) {
      requestBody.email = document.getElementById('email').value;
      requestBody.firstName = document.getElementById('firstName').value;
      requestBody.lastName = document.getElementById('lastName').value;
      requestBody.bio = document.getElementById('bio').value;
      requestBody.interests = document.getElementById('interests').value;
      requestBody.location = document.getElementById('location').value;
      requestBody.latitude = parseFloat(document.getElementById('latitude').value);
      requestBody.longitude = parseFloat(document.getElementById('longitude').value);
      requestBody.age = parseInt(document.getElementById('age').value);
      requestBody.gender = document.getElementById('gender').value;
      requestBody.instagramUsername = document.getElementById('instagramUsername').value;
      requestBody.facebookUsername = document.getElementById('facebookUsername').value;
      requestBody.snapchatUsername = document.getElementById('snapchatUsername').value;
    }

    return requestBody;
  }

  getEndpoint() {
    return this.isSignupForm() ? '/users' : '/login';
  }

  getSuccessMessage() {
    return this.isSignupForm() ? 'Sign up successful' : 'Login successful. Token: ';
  }

  sendRequest(requestBody, endpoint, successMessage) {
    const xhr = new XMLHttpRequest();
    xhr.open('POST', 'http://localhost:8000' + endpoint);
    xhr.setRequestHeader('Content-Type', 'application/json');

    xhr.onload = () => {
      if (xhr.status === 200) {
        const data = JSON.parse(xhr.responseText);
        const token = data.token;
        this.showSuccessMessage(successMessage, token);
        this.storeToken(token);
        this.redirectToMainPage();
      } else {
        this.showErrorMessage('Request failed. Please check your input.');
      }
    };

    xhr.onerror = () => {
      this.showErrorMessage('Request failed. Please check your input.');
    };

    xhr.send(JSON.stringify(requestBody));
  }

  showSuccessMessage(message, token) {
    this.messageElement.textContent = message;
    this.messageElement.style.color = 'green';
  }

  showErrorMessage(message) {
    this.messageElement.textContent = message;
    this.messageElement.style.color = 'red';
  }

  storeToken(token) {
    localStorage.setItem('token', token);
  }

  redirectToMainPage() {
    if (this.isLoginForm()) {
      window.location.href = '/'; // Redirect to events page after successful login
    } else {
      window.location.href = '/login.html'; // Redirect to login page after successful sign up
    }
  }
}

// Initialize the login form
const loginForm = new LoginForm();