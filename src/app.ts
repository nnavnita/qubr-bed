const express =  require("express");
const app = express();
const port = 8080; // default port to listen
const cors = require('cors');
const bodyParser = require('body-parser');
const morgan = require('morgan');
const _ = require('lodash');
const multer = require('multer');
const upload = multer({
  dest: 'uploads/' // this saves your file into a directory called "uploads"
}); 
const cloudinary = require('cloudinary');
require('dotenv').config();

app.use(cors());
app.use(bodyParser.json());
app.use(bodyParser.urlencoded({extended: true}));
app.use(morgan('dev'));

// define a route handler for the default home page
app.post( "/uploadFile", upload.single('avatar'), ( _: any, res: any ) => {
    res.send('file uploaded');
});

cloudinary.config({ 
  cloud_name:process.env.CLOUDINARY_BUCKET_NAME, 
  api_key: process.env.CLOUDINARY_API_KEY, 
  api_secret: process.env.CLOUDINARY_API_SECRET
});

app.post( "/saveToCloudinary", (req: any, res: any) => {
    cloudinary.v2.uploader.upload("uploads/e466bdbd21bacd0664186891753aec86.png",
  { public_id: "icon" }, 
  function(error: any, result: any) {
      if (error) {
          console.error("failed to upload");
          res.send(error);
      }
      console.log(result);
      res.send("uploaded to Cloudinary");
    });
})

app.get( "/diagnostic", ( _: any, res: any ) => {
    res.send( "OK!" );
} );

// start the Express server
app.listen( port, () => {
    // ts
    console.log( `server started at http://localhost:${ port }` );
} );