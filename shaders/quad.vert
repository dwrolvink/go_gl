#version 330 core
layout (location=0) in vec3 aPos;
layout (location=1) in vec2 aTexCoord;

uniform float x;

out vec2 TexCoord;
void main() {
    gl_Position = vec4(aPos.x + x, aPos.y + x, aPos.z, 1.0);
    TexCoord = aTexCoord;
}