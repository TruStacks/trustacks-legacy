import React from 'react';
import {Redirect} from '@docusaurus/router';

import '@fontsource/montserrat'

export default function Home(): JSX.Element {
  return <Redirect to="/overview"></Redirect>
}