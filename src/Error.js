import React from 'react';
import Tabs from "react-bulma-components/lib/components/tabs";

class Error extends React.Component {
    constructor(props) {
        super(props)
        this.title = this.title.bind(this)
    }

    title() {
        if (this.props.error === undefined) {
            return ""
        }
        if (this.props.error.length > 0 ) {
            if (this.props.error[0]["Key"] === "unknown") {
                return "🤔"
            }
            return "‼️"
        }
        return "✅"
    }


    render() {
        return (
            <Tabs.Tab active={this.props.active} onClick={this.props.clickHandler}>{this.props.version} {this.title()}</Tabs.Tab>
        )
    }
}

export default Error
