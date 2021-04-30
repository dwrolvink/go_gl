#version 450
uniform float x;
in vec3 vp;
void main() {
    gl_Position = vec4(vp[0]+x, vp[1]+x, vp[2], 1.0);
}