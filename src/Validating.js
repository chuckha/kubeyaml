import React from 'react';
import Tag from "react-bulma-components/lib/components/tag";

class Validating extends React.Component {
    render() {
        const className = this.props.validating ? "" : "is-hidden"
        return (
            <Tag pull="right" className={className + " mx-0 my-0 px-0 py-0"}>Validating...</Tag>
        )
    }
}

export default Validating