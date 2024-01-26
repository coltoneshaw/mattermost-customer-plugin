import React from 'react';

type ItemParams = {
    title: string;
    value: string;
    editing?: boolean;
    editComponent?: React.ReactNode;
    type: 'string' | 'number';
}
const Item = ({
    title,
    value,
    editing,
    editComponent,
    type,
}: ItemParams) => {
    return (
        <div
            key={title}
            style={{
                display: 'flex',
                flexDirection: 'row',
                gap: '1em',

                // justifyContent: 'space-between',
                alignItems: 'center',
                padding: '4px 0',
                width: '100%',
            }}
        >
            <span
                style={{
                    fontWeight: 600,
                    flex: 1,
                }}
            >
                {title + ':'}
            </span>
            {
                editing ? (
                    <div style={{flex: 2}}>{editComponent}</div>
                ) : (
                    <span>{(type === 'string') ? value : Number(value).toLocaleString() }</span>
                )
            }
        </div>
    );
};

export {
    Item,
};