"use strict";(self.webpackChunk=self.webpackChunk||[]).push([[478],{3905:(t,e,n)=>{n.d(e,{Zo:()=>c,kt:()=>d});var a=n(7294);function o(t,e,n){return e in t?Object.defineProperty(t,e,{value:n,enumerable:!0,configurable:!0,writable:!0}):t[e]=n,t}function r(t,e){var n=Object.keys(t);if(Object.getOwnPropertySymbols){var a=Object.getOwnPropertySymbols(t);e&&(a=a.filter((function(e){return Object.getOwnPropertyDescriptor(t,e).enumerable}))),n.push.apply(n,a)}return n}function i(t){for(var e=1;e<arguments.length;e++){var n=null!=arguments[e]?arguments[e]:{};e%2?r(Object(n),!0).forEach((function(e){o(t,e,n[e])})):Object.getOwnPropertyDescriptors?Object.defineProperties(t,Object.getOwnPropertyDescriptors(n)):r(Object(n)).forEach((function(e){Object.defineProperty(t,e,Object.getOwnPropertyDescriptor(n,e))}))}return t}function l(t,e){if(null==t)return{};var n,a,o=function(t,e){if(null==t)return{};var n,a,o={},r=Object.keys(t);for(a=0;a<r.length;a++)n=r[a],e.indexOf(n)>=0||(o[n]=t[n]);return o}(t,e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(t);for(a=0;a<r.length;a++)n=r[a],e.indexOf(n)>=0||Object.prototype.propertyIsEnumerable.call(t,n)&&(o[n]=t[n])}return o}var p=a.createContext({}),s=function(t){var e=a.useContext(p),n=e;return t&&(n="function"==typeof t?t(e):i(i({},e),t)),n},c=function(t){var e=s(t.components);return a.createElement(p.Provider,{value:e},t.children)},u={inlineCode:"code",wrapper:function(t){var e=t.children;return a.createElement(a.Fragment,{},e)}},m=a.forwardRef((function(t,e){var n=t.components,o=t.mdxType,r=t.originalType,p=t.parentName,c=l(t,["components","mdxType","originalType","parentName"]),m=s(n),d=o,k=m["".concat(p,".").concat(d)]||m[d]||u[d]||r;return n?a.createElement(k,i(i({ref:e},c),{},{components:n})):a.createElement(k,i({ref:e},c))}));function d(t,e){var n=arguments,o=e&&e.mdxType;if("string"==typeof t||o){var r=n.length,i=new Array(r);i[0]=m;var l={};for(var p in e)hasOwnProperty.call(e,p)&&(l[p]=e[p]);l.originalType=t,l.mdxType="string"==typeof t?t:o,i[1]=l;for(var s=2;s<r;s++)i[s]=n[s];return a.createElement.apply(null,i)}return a.createElement.apply(null,n)}m.displayName="MDXCreateElement"},4227:(t,e,n)=>{n.r(e),n.d(e,{assets:()=>p,contentTitle:()=>i,default:()=>u,frontMatter:()=>r,metadata:()=>l,toc:()=>s});var a=n(7462),o=(n(7294),n(3905));const r={title:"Application Creation",slug:"/tutorial/create-application",sidebar_position:3},i="Applications",l={unversionedId:"getting-started/tutorial/create-application",id:"getting-started/tutorial/create-application",title:"Application Creation",description:"TruStacks applications deploy CI/CD assets that consume the software toolchain.",source:"@site/docs/getting-started/tutorial/create-application.md",sourceDirName:"getting-started/tutorial",slug:"/tutorial/create-application",permalink:"/tutorial/create-application",draft:!1,tags:[],version:"current",sidebarPosition:3,frontMatter:{title:"Application Creation",slug:"/tutorial/create-application",sidebar_position:3},sidebar:"defaultSidebar",previous:{title:"Toolchain Installation",permalink:"/tutorial/install-toolchain"},next:{title:"Application Build & Release",permalink:"/tutorial/application build"}},p={},s=[{value:"Container Registry",id:"container-registry",level:2},{value:"Docker Hub",id:"docker-hub",level:3},{value:"Repository host",id:"repository-host",level:4},{value:'<a name="configuration"></a> Configuration',id:"-configuration",level:2},{value:"Workflow Inputs",id:"workflow-inputs",level:3},{value:"Creating the application",id:"creating-the-application",level:3},{value:"Application components",id:"application-components",level:2},{value:"Day 1 Automation",id:"day-1-automation",level:3}],c={toc:s};function u(t){let{components:e,...n}=t;return(0,o.kt)("wrapper",(0,a.Z)({},c,n,{components:e,mdxType:"MDXLayout"}),(0,o.kt)("h1",{id:"applications"},"Applications"),(0,o.kt)("p",null,"TruStacks applications deploy CI/CD assets that consume the software toolchain."),(0,o.kt)("h2",{id:"container-registry"},"Container Registry"),(0,o.kt)("p",null,"The react worfklow for the react-tutorial project requires a container registry."),(0,o.kt)("p",null,"If you already have a docker registry then you can skip to the ",(0,o.kt)("a",{parentName:"p",href:"#configuration"},"Configuration")," section."),(0,o.kt)("h3",{id:"docker-hub"},"Docker Hub"),(0,o.kt)("p",null,"If you do not have a container registry you can use ",(0,o.kt)("a",{parentName:"p",href:"https://hub.docker.com/"},"docker hub"),"."),(0,o.kt)("p",null,"After creating your account follow ",(0,o.kt)("a",{parentName:"p",href:"https://docs.docker.com/docker-hub/access-tokens/"},"this guide")," to create an access token for use in the next step."),(0,o.kt)("h4",{id:"repository-host"},"Repository host"),(0,o.kt)("p",null,"Certain container registries, such as Azure Container Registry, allow dynamic creation of container repositories. "),(0,o.kt)("p",null,"Docker Hub and other container registries requires that repositories be created before pushing images."),(0,o.kt)("p",null,"The format for container images build by the TruStacks CI workflow is:"),(0,o.kt)("pre",null,(0,o.kt)("code",{parentName:"pre"},"<registry-host>/<project-name>:<image-tag>\n")),(0,o.kt)("p",null,"If you are using Docker Hub or a container that requires the repository to exist before pushing images then create the repository now."),(0,o.kt)("h2",{id:"-configuration"},(0,o.kt)("a",{name:"configuration"})," Configuration"),(0,o.kt)("p",null,"To get started, add the following to your ",(0,o.kt)("inlineCode",{parentName:"p"},"react-tutorial-config.yaml")," configuration file:"),(0,o.kt)("pre",null,(0,o.kt)("code",{parentName:"pre",className:"language-yaml"},'applications:\n- name: react-tutorial\n  source: https://github.com/trustacks/workflows\n  workflow: react\n  vars:\n    image: "quay.io/trustacks/workflows"\n    imageTag: "1.0.0"\n    workflow: "react"\n    project: "react-tutorial"\n    gitRemote: "<your fork\'s ssh url>"\n    registryHost: "<your registry hostname>"\n    registryUsername: "<your registry username>"\n  secrets:\n    gitPrivateKey: |- \n      "<your ssh or deploy key>"\n    registryPassword: "<your container registry password or access token>"\n')),(0,o.kt)("p",null,"Let's break down the configuration values:"),(0,o.kt)("blockquote",null,(0,o.kt)("p",{parentName:"blockquote"},(0,o.kt)("inlineCode",{parentName:"p"},"applications")," is the list of available applications to create.",(0,o.kt)("br",{parentName:"p"}),"\n",(0,o.kt)("inlineCode",{parentName:"p"},"applications[*].name")," is the name of the application.",(0,o.kt)("br",{parentName:"p"}),"\n",(0,o.kt)("inlineCode",{parentName:"p"},"applications[*].source")," is the the git repository that contains the workflow resources.",(0,o.kt)("br",{parentName:"p"}),"\n",(0,o.kt)("inlineCode",{parentName:"p"},"applications[*].workflow")," is the workflow to use from the workflow source.",(0,o.kt)("br",{parentName:"p"}),"\n",(0,o.kt)("inlineCode",{parentName:"p"},"applications[*].vars")," are plaintext values that are used by the workflow CI/CD build.",(0,o.kt)("br",{parentName:"p"}),"\n",(0,o.kt)("inlineCode",{parentName:"p"},"applications[*].secrets")," are secret values that are used by the workflow CI/CD build.")),(0,o.kt)("admonition",{type:"tip"},(0,o.kt)("p",{parentName:"admonition"},"Remember to append your project name to ",(0,o.kt)("inlineCode",{parentName:"p"},"registryHost")," if your container registry does not support dynamically creating container repositories."),(0,o.kt)("p",{parentName:"admonition"},(0,o.kt)("em",{parentName:"p"},"ie. ",(0,o.kt)("inlineCode",{parentName:"em"},"registry.hub.docker.com/<project-name>")))),(0,o.kt)("admonition",{type:"caution"},(0,o.kt)("p",{parentName:"admonition"},"Use string quotes to ensure ",(0,o.kt)("inlineCode",{parentName:"p"},"vars")," and ",(0,o.kt)("inlineCode",{parentName:"p"},"secrets")," are interpreted as strings. Numerical values will result in errors.")),(0,o.kt)("admonition",{type:"caution"},(0,o.kt)("p",{parentName:"admonition"},"Your repository url must use ssh, and you must use an ssh key.")),(0,o.kt)("h3",{id:"workflow-inputs"},"Workflow Inputs"),(0,o.kt)("p",null,"Inputs are provided throught ",(0,o.kt)("strong",{parentName:"p"},"vars")," and ",(0,o.kt)("strong",{parentName:"p"},"secrets"),". the input configuration values are passed directly to the CI/CD build."),(0,o.kt)("h3",{id:"creating-the-application"},"Creating the application"),(0,o.kt)("p",null,"After adding the application configuration to ",(0,o.kt)("inlineCode",{parentName:"p"},"react-tutorial-config.yaml"),", we are now ready to create the application."),(0,o.kt)("p",null,"Use the following command create the application: "),(0,o.kt)("pre",null,(0,o.kt)("code",{parentName:"pre",className:"language-bash"},"tsctl application create --name react-tutorial --config react-tutorial-config.yaml\n")),(0,o.kt)("h2",{id:"application-components"},"Application components"),(0,o.kt)("p",null,"In addition to the ",(0,o.kt)("a",{parentName:"p",href:"/tutorial/install-toolchain#-core-components-at-a-glance"},"Core Components")," installed in the toolchain, the ",(0,o.kt)("a",{parentName:"p",href:"https://github.com/TruStacks/workflows/tree/main/workflows/react"},"react")," workflow installs ",(0,o.kt)("a",{parentName:"p",href:"https://argoproj.github.io/cd/"},"Argo CD")," as a supplemental component."),(0,o.kt)("p",null,"Argo CD will be used to deploy the application"),(0,o.kt)("h3",{id:"day-1-automation"},"Day 1 Automation"),(0,o.kt)("p",null,"In addition to the user provided workflow inputs, toolchain components can provide system inputs."),(0,o.kt)("p",null,"Components built for the TruStacks eco-system implement Day 1 automation to make consumption seamless."),(0,o.kt)("p",null,"The Argo CD ",(0,o.kt)("a",{parentName:"p",href:"https://github.com/TruStacks/catalog/blob/main/pkg/components/argocd/component.go"},"component")," creates a service account, confingures rbac, and exposes its server endpoint and service account secret during installation."),(0,o.kt)("p",null,'Activites such as creating SSO client secrets, configuring rbac groups, and all other activities related to the controlled component and "zero touch" consumption readiness are completed during installation. '),(0,o.kt)("admonition",{type:"info"},(0,o.kt)("p",{parentName:"admonition"},"Day 1 Automation is covered in greater detail in the Core Concepts.")))}u.isMDXComponent=!0}}]);