#version 330 core
layout (location=0) in vec3 aPos;
layout (location=1) in vec2 aTexCoord;

uniform float x;
uniform float t;


out vec2 TexCoord;
void main() {
    gl_Position = vec4(aPos.x, aPos.y, aPos.z, 0.2*cos((0.5*t)) + 1.0 );
    TexCoord = aTexCoord;
}