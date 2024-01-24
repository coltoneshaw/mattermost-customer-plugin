import React from 'react';
import {TextInput} from '@mantine/core';

import {FormTextInputParams} from './types';
import {formTextInputStyles} from './styles';

export const FormTextInput = <T, >({
    getInputProps,
    label,
    placeholder,
    formKey,
}: FormTextInputParams<T>) => {
    const {classes} = formTextInputStyles();
    return (
        <TextInput
            bg={'transparent'}
            classNames={classes}
            placeholder={placeholder || label}
            label={label}
            {...getInputProps(formKey)}
        />
    );
};
