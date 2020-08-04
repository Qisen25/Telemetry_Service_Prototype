const passport = require('passport');
const GoogleAuth = require('passport-google-oauth2').Strategy;

passport.use(new GoogleAuth({
	clientID:,
	clientSecret:,
	callbackURL:
});