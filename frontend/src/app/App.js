import React from "react";
import { BrowserRouter } from "react-router-dom";
import Header from './_components/header'
import AppRouter from './app-router'
import Container from 'react-bootstrap/Container'
import 'firebase/analytics';
import {
  FirebaseAppProvider,
} from "reactfire";

const firebaseConfig = {
  apiKey: "AIzaSyByFxIrUyON4YVhqxyJMIn8g2qYpG_ytvc",
  authDomain: "aioz-blockchain.firebaseapp.com",
  projectId: "aioz-blockchain",
  storageBucket: "aioz-blockchain.appspot.com",
  messagingSenderId: "596368929019",
  appId: "1:596368929019:web:ccb58732e2d48036f6c05e",
  measurementId: "G-JDQHGFN4QQ"
};

function App() {
  return (
    <FirebaseAppProvider firebaseConfig={firebaseConfig}>
      <BrowserRouter>
        <Header />
        <Container className='pb-4'>
          <AppRouter />
        </Container>
      </BrowserRouter>
    </FirebaseAppProvider>
  );
}

export default App;
