#version 330 core
layout (location=0) in vec3 aPos;
layout (location=1) in vec2 aTexCoord;

uniform float x;
uniform float y;
uniform float scale;
uniform float t;


out vec2 TexCoord;
void main() {
    gl_Position = vec4(
        aPos.x + x*(1.-1./scale), 
        aPos.y + y*(1.-1./scale), 
        aPos.z, 
        1.0/scale);
    TexCoord = aTexCoord;
}