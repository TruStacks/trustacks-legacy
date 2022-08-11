"use strict";(self.webpackChunk=self.webpackChunk||[]).push([[552],{3905:(t,e,n)=>{n.d(e,{Zo:()=>s,kt:()=>h});var o=n(7294);function a(t,e,n){return e in t?Object.defineProperty(t,e,{value:n,enumerable:!0,configurable:!0,writable:!0}):t[e]=n,t}function i(t,e){var n=Object.keys(t);if(Object.getOwnPropertySymbols){var o=Object.getOwnPropertySymbols(t);e&&(o=o.filter((function(e){return Object.getOwnPropertyDescriptor(t,e).enumerable}))),n.push.apply(n,o)}return n}function r(t){for(var e=1;e<arguments.length;e++){var n=null!=arguments[e]?arguments[e]:{};e%2?i(Object(n),!0).forEach((function(e){a(t,e,n[e])})):Object.getOwnPropertyDescriptors?Object.defineProperties(t,Object.getOwnPropertyDescriptors(n)):i(Object(n)).forEach((function(e){Object.defineProperty(t,e,Object.getOwnPropertyDescriptor(n,e))}))}return t}function l(t,e){if(null==t)return{};var n,o,a=function(t,e){if(null==t)return{};var n,o,a={},i=Object.keys(t);for(o=0;o<i.length;o++)n=i[o],e.indexOf(n)>=0||(a[n]=t[n]);return a}(t,e);if(Object.getOwnPropertySymbols){var i=Object.getOwnPropertySymbols(t);for(o=0;o<i.length;o++)n=i[o],e.indexOf(n)>=0||Object.prototype.propertyIsEnumerable.call(t,n)&&(a[n]=t[n])}return a}var c=o.createContext({}),p=function(t){var e=o.useContext(c),n=e;return t&&(n="function"==typeof t?t(e):r(r({},e),t)),n},s=function(t){var e=p(t.components);return o.createElement(c.Provider,{value:e},t.children)},u={inlineCode:"code",wrapper:function(t){var e=t.children;return o.createElement(o.Fragment,{},e)}},d=o.forwardRef((function(t,e){var n=t.components,a=t.mdxType,i=t.originalType,c=t.parentName,s=l(t,["components","mdxType","originalType","parentName"]),d=p(n),h=a,m=d["".concat(c,".").concat(h)]||d[h]||u[h]||i;return n?o.createElement(m,r(r({ref:e},s),{},{components:n})):o.createElement(m,r({ref:e},s))}));function h(t,e){var n=arguments,a=e&&e.mdxType;if("string"==typeof t||a){var i=n.length,r=new Array(i);r[0]=d;var l={};for(var c in e)hasOwnProperty.call(e,c)&&(l[c]=e[c]);l.originalType=t,l.mdxType="string"==typeof t?t:a,r[1]=l;for(var p=2;p<i;p++)r[p]=n[p];return o.createElement.apply(null,r)}return o.createElement.apply(null,n)}d.displayName="MDXCreateElement"},8149:(t,e,n)=>{n.r(e),n.d(e,{assets:()=>c,contentTitle:()=>r,default:()=>u,frontMatter:()=>i,metadata:()=>l,toc:()=>p});var o=n(7462),a=(n(7294),n(3905));const i={title:"Toolchain Installation",slug:"/tutorial/install-toolchain",sidebar_position:2},r="Software Toolchain",l={unversionedId:"getting-started/tutorial/install-toolchain",id:"getting-started/tutorial/install-toolchain",title:"Toolchain Installation",description:"The toolchain contains the software components required to build and release applications.",source:"@site/docs/getting-started/tutorial/install-toolchain.md",sourceDirName:"getting-started/tutorial",slug:"/tutorial/install-toolchain",permalink:"/tutorial/install-toolchain",draft:!1,tags:[],version:"current",sidebarPosition:2,frontMatter:{title:"Toolchain Installation",slug:"/tutorial/install-toolchain",sidebar_position:2},sidebar:"defaultSidebar",previous:{title:"Project Setup",permalink:"/tutorial/project-setup"},next:{title:"Application Creation",permalink:"/tutorial/create-application"}},c={},p=[{value:"Configuration",id:"configuration",level:2},{value:'<a id="core-components-aag"></a> Core Components At-A-Glance',id:"-core-components-at-a-glance",level:2},{value:"Install",id:"install",level:2},{value:"SSO Setup",id:"sso-setup",level:2},{value:"CI Provider",id:"ci-provider",level:2}],s={toc:p};function u(t){let{components:e,...n}=t;return(0,a.kt)("wrapper",(0,o.Z)({},s,n,{components:e,mdxType:"MDXLayout"}),(0,a.kt)("h1",{id:"software-toolchain"},"Software Toolchain"),(0,a.kt)("p",null,"The toolchain contains the software components required to build and release applications."),(0,a.kt)("h2",{id:"configuration"},"Configuration"),(0,a.kt)("p",null,"To install a new toolchain we need to create a toolchain configuration."),(0,a.kt)("p",null,"Create a file named ",(0,a.kt)("inlineCode",{parentName:"p"},"react-tutorial-config.yaml")," and add the following values:"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre"},"name: react-tutorial\nsource: https://github.com/trustacks/toolchain\nparameters:\n  sso: authentik\n  ingressPort: 8081 # change this value if you used a different port for your k3d loadbalaner\n  tls: false\n  ci: concourse\n")),(0,a.kt)("p",null,"Let's break down the configuration values:"),(0,a.kt)("blockquote",null,(0,a.kt)("p",{parentName:"blockquote"},(0,a.kt)("inlineCode",{parentName:"p"},"name")," is the name of the toolchain.",(0,a.kt)("br",{parentName:"p"}),"\n",(0,a.kt)("inlineCode",{parentName:"p"},"source")," is the the git repository that contains the toolchain resources.",(0,a.kt)("br",{parentName:"p"}),"\n",(0,a.kt)("inlineCode",{parentName:"p"},"parameters")," are ",(0,a.kt)("a",{parentName:"p",href:"https://github.com/TruStacks/catalog/blob/main/pkg/catalog/catalog.yaml"},"values")," that will be passed to the software components during installation.")),(0,a.kt)("admonition",{title:"in case you missed it",type:"caution"},(0,a.kt)("p",{parentName:"admonition"},"If you changed your ingress port during the ",(0,a.kt)("a",{parentName:"p",href:"http://localhost:3000/installation#cluster-creation"},"k3d cluster creation")," to something other than ",(0,a.kt)("inlineCode",{parentName:"p"},"8081"),", make sure to update ",(0,a.kt)("inlineCode",{parentName:"p"},"parameters.ingressPort")," in the yaml configuration file before proceeding.")),(0,a.kt)("h2",{id:"-core-components-at-a-glance"},(0,a.kt)("a",{id:"core-components-aag"})," Core Components At-A-Glance"),(0,a.kt)("p",null,"TruStacks toolchains must include a set of core components as dependencies."),(0,a.kt)("p",null,"The core componets are:"),(0,a.kt)("ul",null,(0,a.kt)("li",{parentName:"ul"},"CI Provider"),(0,a.kt)("li",{parentName:"ul"},"SSO Provider")),(0,a.kt)("p",null,"Toolchains can include any number of supplemental components that provide a wide range of software delivery capabilities."),(0,a.kt)("p",null,"The toolchain used in this tutorial uses ",(0,a.kt)("a",{parentName:"p",href:"https://goauthentik.io/"},"authentik")," for SSO and ",(0,a.kt)("a",{parentName:"p",href:"https://concourse-ci.org/"},"concourse")," for CI."),(0,a.kt)("admonition",{title:"don't like the tools?",type:"tip"},(0,a.kt)("p",{parentName:"admonition"},"TruStacks Engine is built to be extensible. If it can be deployed with helm, then it can likely be built into the TruStacks eco-system."),(0,a.kt)("p",{parentName:"admonition"},"Drop us a suggestion on ",(0,a.kt)("a",{parentName:"p",href:"https://discord.gg/tgpWURqY"},"discord"),", or check our project board to see ",(0,a.kt)("a",{parentName:"p",href:"https://trello.com/b/IwJMgZiO/trustacks-oss"},"what's happening now"),".")),(0,a.kt)("admonition",{type:"info"},(0,a.kt)("p",{parentName:"admonition"},"Parameters, Core Components and other in-depth topics are covered in greater detail in the Core Concepts.")),(0,a.kt)("h2",{id:"install"},"Install"),(0,a.kt)("p",null,"Now that we have our configuration file we can install our toolchain."),(0,a.kt)("p",null,"Run the following command to install the toolchain:"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre"},"tsctl toolchain install --config react-tutorial-config.yaml\n")),(0,a.kt)("p",null,"Check the status of the services with the following command. Wait until all service are in the ",(0,a.kt)("inlineCode",{parentName:"p"},"Running")," state:"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre"},"kubectl get po -n trustacks-toolchain-react-tutorial  \n")),(0,a.kt)("p",null,"Example output:"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre"},"NAME                                READY   STATUS    RESTARTS   AGE\ndind-88db84bd6-8bxxf                1/1     Running   0          91s\nauthentik-worker-78c5d654c4-jlv45   1/1     Running   0          87s\nauthentik-postgresql-0              1/1     Running   0          87s\nauthentik-redis-master-0            1/1     Running   0          87s\nauthentik-server-786cb79b-xq76s     1/1     Running   0          87s\nconcourse-worker-0                  1/1     Running   0          22s\nconcourse-worker-1                  1/1     Running   0          22s\nconcourse-postgresql-0              1/1     Running   0          22s\nconcourse-web-747c56c56f-b94ql      2/2     Running   0          22s\n")),(0,a.kt)("admonition",{title:"air gapped installation",type:"caution"},(0,a.kt)("p",{parentName:"admonition"},"Air-Gapped environemnts are not currently supported.")),(0,a.kt)("h2",{id:"sso-setup"},"SSO Setup"),(0,a.kt)("p",null,"Navigate to ",(0,a.kt)("a",{parentName:"p",href:"http://authentik.local.gd:8081/if/flow/initial-setup/"},"authentik")," to configure the ",(0,a.kt)("inlineCode",{parentName:"p"},"akadmin")," user."),(0,a.kt)("img",{src:"/img/authentik-initial-setup.jpg"}),(0,a.kt)("p",null,"You should see the following page after navigating to the address. Enter an email and password and click continue."),(0,a.kt)("p",null,"After clicking continue you will be taken to the authentik landing page. The sso provider configuration is now complete."),(0,a.kt)("admonition",{title:"changed port?",type:"tip"},(0,a.kt)("p",{parentName:"admonition"},"Navigate to ",(0,a.kt)("inlineCode",{parentName:"p"},"http://authentik.local.gd:<port>/if/flow/initial-setup/")," if you changed your loadbalancer port.")),(0,a.kt)("h2",{id:"ci-provider"},"CI Provider"),(0,a.kt)("p",null,"Now that your sso provider is configured, you can log in to ",(0,a.kt)("a",{parentName:"p",href:"http://concourse.local.gd:8081"},"concourse"),"."),(0,a.kt)("img",{src:"/img/concourse-login.jpg"}),(0,a.kt)("p",null,"Click the login page at the top right of the page or click the login link in the hero image."),(0,a.kt)("img",{src:"/img/concourse-sso.jpg"}),(0,a.kt)("p",null,"Click ",(0,a.kt)("inlineCode",{parentName:"p"},"sso")," to sign in using the sso provider."),(0,a.kt)("admonition",{title:"changed port?",type:"tip"},(0,a.kt)("p",{parentName:"admonition"},"Navigate to ",(0,a.kt)("inlineCode",{parentName:"p"},"http://concourse.local.gd:<port>/if/flow/initial-setup/")," if you changed your loadbalancer port.")))}u.isMDXComponent=!0}}]);