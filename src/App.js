import React from 'react';
import './App.scss';
import Button from "react-bulma-components/lib/components/button";
import Section from "react-bulma-components/lib/components/section"
import Hero from "react-bulma-components/lib/components/hero";
import Container from "react-bulma-components/lib/components/container";
import Heading from "react-bulma-components/lib/components/heading";
import Content from "react-bulma-components/lib/components/content";
import Icon from "react-bulma-components/lib/components/icon";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faExclamationTriangle} from "@fortawesome/free-solid-svg-icons"
import Tile from "react-bulma-components/lib/components/tile";
import {Textarea} from "react-bulma-components/lib/components/form"
import Tabs from "react-bulma-components/lib/components/tabs";

import Validation from "./validation";
import Error from "./Error";
import Document from "./Document";

class App extends React.Component {
    constructor(props) {
        super(props)
        this.state = {
            document: "",
            versions:  ["1.18", "1.17", "1.16", "1.15"],
            active: "1.18",
            errors: {},
        }
        this.validator = new Validation({baseURL: "http://localhost:9000"})

        this.setExample = this.setExample.bind(this)
        this.onChange = this.onChange.bind(this)
        this.handleValidate = this.handleValidate.bind(this)
        this.handleTabClick = this.handleTabClick.bind(this)
        this.errorsCallback = this.errorsCallback.bind(this)
        this.activeError = this.activeError.bind(this)
        this.docError = this.docError.bind(this)
    }

    handleValidate(e) {
        e.preventDefault()
        this.validator.validate(this.state.document, this.errorsCallback)
    }

    handleTabClick(version) {
        return () => {
            this.setState({active: version})
        }
    }

    errorsCallback(response) {
        this.setState({errors: JSON.parse(response)})
    }

    onChange(event) {
        this.setState({document: event.target.value})
    }

    activeError() {
        const errObj = this.state.errors[this.state.active]
        if (errObj === undefined) {
            return ""
        }
        if (errObj.length === 0) {
            return ""
        }
        if (errObj[0].hasOwnProperty("Error")) {
            return errObj[0]["Error"]
        }
        return ""
    }

    docError() {
        const errObj = this.state.errors[this.state.active]
        if (errObj === undefined) {
            return ""
        }
        if (errObj.length === 0) {
            return ""
        }
        if (errObj[0].hasOwnProperty("Error")) {
            let doc = new Document(this.state.document, errObj[0])
            return doc.doc.toString()
        }
        return ""
    }

    setExample(event) {
        event.preventDefault()
        this.setState({document: `apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
    app: nginx
spec:
  replicas: invalid number
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1.7.9
        ports:
        - containerPort: 80
`})
    }

    render() {
        const errorTabs = this.state.versions.map((version) =>
            <Error error={this.state.errors[version]} active={this.state.active === version} version={version} key={version} clickHandler={this.handleTabClick(version)}/>
        )
        return (
            <div className="App">
                <Section>
                    <Hero backgroundColor="light">
                        <Hero.Body>
                            <Container>
                                <Heading size={1}>
                                    Kube YAML
                                </Heading>
                                <Heading subtitle size={3}>
                                    Validating Kubernetes objects since 2018
                                </Heading>
                            </Container>
                        </Hero.Body>
                    </Hero>
                </Section>
                <Container>
                    <Content>
                        <p className="is-size-7">
                            <Icon color="warning">
                                <FontAwesomeIcon icon={faExclamationTriangle}/>
                            </Icon>
                            Please only enter one YAML document at a time.
                        </p>
                    </Content>
                </Container>
                <Section>
                    <Tile kind="ancestor">
                        <Tile kind="parent" size={5}>
                            <Tile kind="child">
                                <Content>
                                    <form>
                                        <Button backgroundColor="success" onClick={this.setExample}>Example YAML</Button>
                                        <Button backgroundColor="info" pull="right" onClick={this.handleValidate}>Validate</Button>
                                        <Textarea rows={30} className="is-family-code" placeholder="Paste YAML here!" onChange={this.onChange} value={this.state.document} />
                                    </form>
                                </Content>
                            </Tile>
                        </Tile>
                        <Tile kind="parent">
                            <Tile kind="child">
                                <Tabs fullwidth>
                                    {errorTabs}
                                </Tabs>
                                <Content>
                                    <p>{this.activeError()}</p>
                                    <pre dangerouslySetInnerHTML={{__html: this.docError()}} />
                                </Content>
                            </Tile>
                        </Tile>
                    </Tile>
                </Section>
            </div>
        );
    }
}

export default App;
